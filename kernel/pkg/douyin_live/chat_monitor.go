package douyin_live

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/elliotchance/orderedmap"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"google.golang.org/protobuf/proto"
	"io"
	"kernel/pkg/douyin_live/generated/douyin"
	"kernel/pkg/douyin_live/js"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func emitProcess[T proto.Message](logger *zap.Logger, message *douyin.Message, msg T, emit func(msg T)) {
	if err := proto.Unmarshal(message.Payload, msg); err != nil {
		logger.Warn("unmarshal fail message", zap.Error(err), zap.String("method", message.Method))
	}
	emit(msg)
}

var (
	rootIdRegexp       = regexp.MustCompile(`roomId\\":\\"(\d+)\\"`)
	userUniqueIdRegexp = regexp.MustCompile(`user_unique_id\\":\\"(\d+)\\"`)
)

func generateMsToken(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+="

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		// 生成0到charset长度之间的随机数
		randInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))

		// 将随机数转换为字符集中的字符
		b[i] = charset[randInt.Int64()]
	}

	return string(b) + "=_"
}

type ChatMonitor struct {
	logger        *zap.Logger
	client        *resty.Client
	ttwId         string
	roomId        string
	liveId        string
	liveUrl       string
	Useragent     string
	userUniqueId  string
	Conn          *websocket.Conn
	eventHandlers []*ChatHandler
}

func (m *ChatMonitor) init() error {
	response, err := m.client.R().Get(m.liveUrl)
	if err != nil {
		m.logger.Error("failed to init live", zap.Error(err))
		return err
	}
	for _, cookie := range response.Cookies() {
		if cookie.Name == "ttwid" {
			m.logger.Debug("ttwId", zap.String("ttwId", cookie.Value))
			m.ttwId = cookie.Value
			break
		}
	}

	//cookie := "ttwid=" + t.Value + "&msToken=" + utils.GenerateMsToken(107) + "; __ac_nonce=0123407cc00a9e438deb4"
	res, err := m.client.R().SetCookies([]*http.Cookie{
		{
			Name:  "ttwid",
			Value: "ttwid=" + m.ttwId + "&msToken=" + generateMsToken(107),
		},
		{
			Name:  "__ac_nonce",
			Value: "0123407cc00a9e438deb4",
		},
	}).Get(m.liveUrl + m.liveId)
	if err != nil {
		m.logger.Error("failed to init live", zap.Error(err))
		return err

	}
	rootIdMatch := rootIdRegexp.FindStringSubmatch(res.String())
	m.roomId = rootIdMatch[1]
	userUniqueIdMatch := userUniqueIdRegexp.FindStringSubmatch(res.String())
	m.userUniqueId = userUniqueIdMatch[1]
	return nil
}

func (m *ChatMonitor) monitor() {
	for {
		_, message, err := m.Conn.ReadMessage()
		if err != nil {
			m.logger.Warn("failed to read message, maybe it's off.", zap.Error(err), zap.Binary("message", message))
			return
		}
		if message != nil {
			pac := &douyin.PushFrame{}
			err := proto.Unmarshal(message, pac)
			if err != nil {
				m.logger.Warn("failed to unmarshal push frame", zap.Error(err), zap.Binary("message", message))
				continue
			}
			n := false
			for _, v := range pac.HeadersList {
				if v.Key == "compress_type" {
					if v.Value == "gzip" {
						n = true
						continue
					}
				}
			}
			//消息为gzip压缩进行处理.否则抛弃
			if n == true && pac.PayloadType == "msg" {
				gzipReader, err := gzip.NewReader(bytes.NewReader(pac.Payload))
				if err != nil {
					m.logger.Warn("failed to create gzip reader", zap.Error(err), zap.Binary("payload", pac.Payload))
					continue
				}
				uncompressedData, err := io.ReadAll(gzipReader)
				gzipReader.Close()
				if err != nil {
					m.logger.Warn("failed to read gzip", zap.Error(err), zap.Binary("payload", pac.Payload))
					continue
				}
				response := &douyin.Response{}
				err = proto.Unmarshal(uncompressedData, response)
				if err != nil {
					m.logger.Warn("failed to unmarshal response", zap.Error(err), zap.Binary("payload", uncompressedData))
					continue
				}
				if response.NeedAck {
					ack := &douyin.PushFrame{
						LogId:       pac.LogId,
						PayloadType: "ack",
						Payload:     []byte(response.InternalExt),
					}
					serializedAck, err := proto.Marshal(ack)
					if err != nil {
						m.logger.Warn("unable to serialize heartbeat packets", zap.Error(err))
					}
					err = m.Conn.WriteMessage(websocket.BinaryMessage, serializedAck)
					if err != nil {
						m.logger.Warn("failed to sending heartbeat packet", zap.Error(err))
					}
				}

				m.emit(response)
			}

		}
	}
}

func (m *ChatMonitor) emit(response *douyin.Response) {
	for _, message := range response.MessagesList {
		switch message.Method {
		case "WebcastChatMessage":
			emitProcess(m.logger, message, &douyin.ChatMessage{}, func(msg *douyin.ChatMessage) {
				for _, handler := range m.eventHandlers {
					handler.chatHandler(msg)
				}
			})
		case "WebcastGiftMessage":
			emitProcess(m.logger, message, &douyin.GiftMessage{}, func(msg *douyin.GiftMessage) {
				for _, handler := range m.eventHandlers {
					handler.giftHandler(msg)
				}
			})
		case "WebcastLikeMessage":
			emitProcess(m.logger, message, &douyin.LikeMessage{}, func(msg *douyin.LikeMessage) {
				for _, handler := range m.eventHandlers {
					handler.likeHandler(msg)
				}
			})

		case "WebcastMemberMessage":
			emitProcess(m.logger, message, &douyin.MemberMessage{}, func(msg *douyin.MemberMessage) {
				for _, handler := range m.eventHandlers {
					handler.memberHandler(msg)
				}
			})

		case "WebcastSocialMessage":
			emitProcess(m.logger, message, &douyin.SocialMessage{}, func(msg *douyin.SocialMessage) {
				for _, handler := range m.eventHandlers {
					handler.socialHandler(msg)
				}
			})

		case "WebcastRoomUserSeqMessage":
			emitProcess(m.logger, message, &douyin.RoomUserSeqMessage{}, func(msg *douyin.RoomUserSeqMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomUserSeqHandler(msg)
				}
			})

		case "WebcastFansclubMessage":
			emitProcess(m.logger, message, &douyin.FansclubMessage{}, func(msg *douyin.FansclubMessage) {
				for _, handler := range m.eventHandlers {
					handler.fansclubHandler(msg)
				}
			})

		case "WebcastControlMessage":
			emitProcess(m.logger, message, &douyin.ControlMessage{}, func(msg *douyin.ControlMessage) {
				for _, handler := range m.eventHandlers {
					handler.controlHandler(msg)
				}
			})

		case "WebcastEmojiChatMessage":
			emitProcess(m.logger, message, &douyin.EmojiChatMessage{}, func(msg *douyin.EmojiChatMessage) {
				for _, handler := range m.eventHandlers {
					handler.emojiChatHandler(msg)
				}
			})

		case "WebcastRoomStatsMessage":
			emitProcess(m.logger, message, &douyin.RoomStatsMessage{}, func(msg *douyin.RoomStatsMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomStatsHandler(msg)
				}
			})

		case "WebcastRoomMessage":
			emitProcess(m.logger, message, &douyin.RoomMessage{}, func(msg *douyin.RoomMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomHandler(msg)
				}
			})

		case "WebcastRanklistHourEntranceMessage":
			emitProcess(m.logger, message, &douyin.RanklistHourEntranceMessage{}, func(msg *douyin.RanklistHourEntranceMessage) {
				for _, handler := range m.eventHandlers {
					handler.ranklistHourEntranceHandler(msg)
				}
			})

		case "WebcastRoomRankMessage":
			emitProcess(m.logger, message, &douyin.RoomRankMessage{}, func(msg *douyin.RoomRankMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomRankHandler(msg)
				}
			})

		case "WebcastInRoomBannerMessage":
			emitProcess(m.logger, message, &douyin.InRoomBannerMessage{}, func(msg *douyin.InRoomBannerMessage) {
				for _, handler := range m.eventHandlers {
					handler.inRoomBannerHandler(msg)
				}
			})

		case "WebcastRoomDataSyncMessage":
			emitProcess(m.logger, message, &douyin.RoomDataSyncMessage{}, func(msg *douyin.RoomDataSyncMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomDataSyncHandler(msg)
				}
			})

		case "WebcastLuckyBoxTempStatusMessage":
			emitProcess(m.logger, message, &douyin.LuckyBoxTempStatusMessage{}, func(msg *douyin.LuckyBoxTempStatusMessage) {
				for _, handler := range m.eventHandlers {
					handler.luckyBoxTempStatusHandler(msg)
				}
			})

		case "WebcastDecorationModifyMethod":
			emitProcess(m.logger, message, &douyin.DecorationModifyMessage{}, func(msg *douyin.DecorationModifyMessage) {
				for _, handler := range m.eventHandlers {
					handler.decorationModifyHandler(msg)
				}
			})

		case "WebcastLinkMicAudienceKtvMessage":
			emitProcess(m.logger, message, &douyin.LinkMicAudienceKtvMessage{}, func(msg *douyin.LinkMicAudienceKtvMessage) {
				for _, handler := range m.eventHandlers {
					handler.linkMicAudienceKtvHandler(msg)
				}
			})

		case "WebcastRoomStreamAdaptationMessage":
			emitProcess(m.logger, message, &douyin.RoomStreamAdaptationMessage{}, func(msg *douyin.RoomStreamAdaptationMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomStreamAdaptationHandler(msg)
				}
			})

		case "WebcastQuizAudienceStatusMessage":
			emitProcess(m.logger, message, &douyin.QuizAudienceStatusMessage{}, func(msg *douyin.QuizAudienceStatusMessage) {
				for _, handler := range m.eventHandlers {
					handler.quizAudienceStatusHandler(msg)
				}
			})

		case "WebcastHotChatMessage":
			emitProcess(m.logger, message, &douyin.HotChatMessage{}, func(msg *douyin.HotChatMessage) {
				for _, handler := range m.eventHandlers {
					handler.hotChatHandler(msg)
				}
			})

		case "WebcastHotRoomMessage":
			emitProcess(m.logger, message, &douyin.HotRoomMessage{}, func(msg *douyin.HotRoomMessage) {
				for _, handler := range m.eventHandlers {
					handler.hotRoomHandler(msg)
				}
			})

		case "WebcastAudioChatMessage":
			emitProcess(m.logger, message, &douyin.AudioChatMessage{}, func(msg *douyin.AudioChatMessage) {
				for _, handler := range m.eventHandlers {
					handler.audioChatHandler(msg)
				}
			})

		case "WebcastRoomNotifyMessage":
			emitProcess(m.logger, message, &douyin.NotifyMessage{}, func(msg *douyin.NotifyMessage) {
				for _, handler := range m.eventHandlers {
					handler.roomNotifyHandler(msg)
				}
			})

		case "WebcastLuckyBoxMessage":
			emitProcess(m.logger, message, &douyin.LuckyBoxMessage{}, func(msg *douyin.LuckyBoxMessage) {
				for _, handler := range m.eventHandlers {
					handler.luckyBoxHandler(msg)
				}
			})

		case "WebcastUpdateFanTicketMessage":
			emitProcess(m.logger, message, &douyin.UpdateFanTicketMessage{}, func(msg *douyin.UpdateFanTicketMessage) {
				for _, handler := range m.eventHandlers {
					handler.updateFanTicketHandler(msg)
				}
			})

		case "WebcastScreenChatMessage":
			emitProcess(m.logger, message, &douyin.ScreenChatMessage{}, func(msg *douyin.ScreenChatMessage) {
				for _, handler := range m.eventHandlers {
					handler.screenChatHandler(msg)
				}
			})

		case "WebcastNotifyEffectMessage":
			emitProcess(m.logger, message, &douyin.NotifyEffectMessage{}, func(msg *douyin.NotifyEffectMessage) {
				for _, handler := range m.eventHandlers {
					handler.notifyEffectHandler(msg)
				}
			})

		case "WebcastBindingGiftMessage":
			emitProcess(m.logger, message, &douyin.NotifyEffectMessage_BindingGiftMessage{}, func(msg *douyin.NotifyEffectMessage_BindingGiftMessage) {
				for _, handler := range m.eventHandlers {
					handler.bindingGiftHandler(msg)
				}
			})

		case "WebcastTempStateAreaReachMessage":
			emitProcess(m.logger, message, &douyin.TempStateAreaReachMessage{}, func(msg *douyin.TempStateAreaReachMessage) {
				for _, handler := range m.eventHandlers {
					handler.tempStateAreaReachHandler(msg)
				}
			})

		case "WebcastGrowthTaskMessage":
			emitProcess(m.logger, message, &douyin.GrowthTaskMessage{}, func(msg *douyin.GrowthTaskMessage) {
				for _, handler := range m.eventHandlers {
					handler.growthTaskHandler(msg)
				}
			})

		case "WebcastGameCPBaseMessage":
			emitProcess(m.logger, message, &douyin.GameCPBaseMessage{}, func(msg *douyin.GameCPBaseMessage) {
				for _, handler := range m.eventHandlers {
					handler.gameCPBaseHandler(msg)
				}
			})

		default:
			for _, handler := range m.eventHandlers {
				handler.unknownHandler(message)
			}
		}
	}
	//method := data.Method
	//
	//switch method {
	//
	//case "WebcastChatMessage":
	//	//msg := &douyin.ChatMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Println("聊天msg", msg.User.Id, msg.User.NickName, msg.Content)
	//case "WebcastGiftMessage":
	//	//msg := &douyin.GiftMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Println("礼物msg", msg.User.Id, msg.User.NickName, msg.Gift.Name, msg.ComboCount)
	//case "WebcastLikeMessage":
	//	//msg := &douyin.LikeMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Println("点赞msg", msg.User.Id, msg.User.NickName, msg.Count)
	//case "WebcastMemberMessage":
	//	//msg := &douyin.MemberMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Println("进场msg", msg.User.Id, msg.User.NickName, msg.User.Gender)
	//case "WebcastSocialMessage":
	//	//msg := &douyin.SocialMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Println("关注msg", msg.User.Id, msg.User.NickName)
	//case "WebcastRoomUserSeqMessage":
	//	//msg := &douyin.RoomUserSeqMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("房间人数msg 当前观看人数:%v,累计观看人数:%v\n", msg.Total, msg.TotalPvForAnchor)
	//case "WebcastFansclubMessage":
	//	//msg := &douyin.FansclubMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("粉丝团msg %v\n", msg.Content)
	//case "WebcastControlMessage":
	//	//msg := &douyin.ControlMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("直播间状态消息%v", msg.Status)
	//case "WebcastEmojiChatMessage":
	//	//msg := &douyin.EmojiChatMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("表情消息%vuser:%vcommon:%vdefault_content:%v", msg.EmojiId, msg.User, msg.Common, msg.DefaultContent)
	//case "WebcastRoomStatsMessage":
	//	//msg := &douyin.RoomStatsMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("直播间统计msg%v", msg.DisplayLong)
	//case "WebcastRoomMessage":
	//	//msg := &douyin.RoomMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("【直播间msg】直播间id%v", msg.Common.RoomId)
	//case "WebcastRoomRankMessage":
	//	//msg := &douyin.RoomRankMessage{}
	//	//proto.Unmarshal(data.Payload, msg)
	//	//log.Printf("直播间排行榜msg%v", msg.RanksList)
	//
	//default:
	//	//m.emit(Default, data.Payload)
	//	//log.Println("payload:", method, hex.EncodeToString(data.Payload))
	//}

}

func (m *ChatMonitor) Start(handlers ...*ChatHandler) {
	var err error

	smap := orderedmap.NewOrderedMap()
	smap.Set("live_id", "1")
	smap.Set("aid", "6383")
	smap.Set("version_code", "180800")
	smap.Set("webcast_sdk_version", "1.0.14-beta.0")
	smap.Set("room_id", m.roomId)
	smap.Set("sub_room_id", "")
	smap.Set("sub_channel_id", "")
	smap.Set("did_rule", "3")
	smap.Set("user_unique_id", m.userUniqueId)
	smap.Set("device_platform", "web")
	smap.Set("device_type", "")
	smap.Set("ac", "")
	smap.Set("identity", "audience")
	signature := js.Signature(smap)

	browserInfo := strings.Split(m.Useragent, "Mozilla")[1]
	parsedURL := strings.Replace(browserInfo[1:], " ", "%20", -1)
	fetchTime := time.Now().UnixNano() / int64(time.Millisecond)

	browserVersion := parsedURL

	wssURL := "wss://webcast5-ws-web-lf.douyin.com/webcast/im/push/v2/?app_name=douyin_web&version_code=180800&" +
		"webcast_sdk_version=1.0.14-beta.0&update_version_code=1.0.14-beta.0&compress=gzip&device_platform" +
		"=web&cookie_enabled=true&screen_width=1920&screen_height=1080&browser_language=zh-CN&browser_platform=Win32&" +
		"browser_name=Mozilla&browser_version=" + browserVersion + "&browser_online=true" +
		"&tz_name=Asia/Shanghai&cursor=m-1_u-1_fh-7383731312643626035_t-1719159695790_r-1&internal_ext" +
		"=internal_src:dim|wss_push_room_id:" + m.roomId + "|wss_push_did:" + m.userUniqueId + "|first_req_ms:" + cast.ToString(fetchTime) + "|fetch_time:" + cast.ToString(fetchTime) + "|seq:1|wss_info:0-" + cast.ToString(fetchTime) + "-0-0|" +
		"wrds_v:7382620942951772256&host=https://live.douyin.com&aid=6383&live_id=1&did_rule=3" +
		"&endpoint=live_pc&support_wrds=1&user_unique_id=" + m.userUniqueId + "&im_path=/webcast/im/fetch/" +
		"&identity=audience&need_persist_msg_count=15&insert_task_id=&live_reason=&room_id=" + m.roomId + "&heartbeatDuration=0&signature=" + signature

	headers := make(http.Header)
	headers.Add("user-agent", m.Useragent)
	headers.Add("cookie", fmt.Sprintf("ttwid=%s", m.ttwId))
	var response *http.Response
	m.Conn, response, err = websocket.DefaultDialer.Dial(wssURL, headers)

	if err != nil {
		m.logger.Error("failed to dial wss", zap.Error(err), zap.String("roomId", m.roomId), zap.String("ttwId", m.ttwId), zap.String("wssURL", wssURL), zap.Int("statusCode", response.StatusCode))
		return
	}
	m.logger.Debug("successful connection")
	m.eventHandlers = handlers
	m.monitor()

}

func New(liveId string) (*ChatMonitor, error) {

	ua := browser.Chrome()

	client := resty.New()
	client.SetHeader("User-Agent", ua)
	if err := js.Init("/Users/klein/Projects/matrix/kernel/pkg/douyin_live/js/webmssdk.js", ua); err != nil {
		return nil, err
	}
	monitor := &ChatMonitor{
		logger:        zap.L().Named(fmt.Sprintf("douyin-monitor-%s", liveId)),
		liveUrl:       "https://live.douyin.com/",
		liveId:        liveId,
		Useragent:     ua,
		client:        client,
		eventHandlers: make([]*ChatHandler, 0),
	}
	return monitor, monitor.init()
}
