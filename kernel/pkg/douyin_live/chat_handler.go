package douyin_live

import "kernel/pkg/douyin_live/generated/douyin"

type EventHandlerOption func(handler *ChatHandler)

type ChatHandler struct {
	_chatHandler                 func(message *douyin.ChatMessage)
	_giftHandler                 func(message *douyin.GiftMessage)
	_likeHandler                 func(message *douyin.LikeMessage)
	_memberHandler               func(message *douyin.MemberMessage)
	_socialHandler               func(message *douyin.SocialMessage)
	_roomUserSeqHandler          func(message *douyin.RoomUserSeqMessage)
	_fansclubHandler             func(message *douyin.FansclubMessage)
	_controlHandler              func(message *douyin.ControlMessage)
	_emojiChatHandler            func(message *douyin.EmojiChatMessage)
	_roomStatsHandler            func(message *douyin.RoomStatsMessage)
	_roomHandler                 func(message *douyin.RoomMessage)
	_ranklistHourEntranceHandler func(message *douyin.RanklistHourEntranceMessage)
	_roomRankHandler             func(message *douyin.RoomRankMessage)
	_inRoomBannerHandler         func(message *douyin.InRoomBannerMessage)
	_roomDataSyncHandler         func(message *douyin.RoomDataSyncMessage)
	_luckyBoxTempStatusHandler   func(message *douyin.LuckyBoxTempStatusMessage)
	_decorationModifyHandler     func(message *douyin.DecorationModifyMessage)
	_linkMicAudienceKtvHandler   func(message *douyin.LinkMicAudienceKtvMessage)
	_roomStreamAdaptationHandler func(message *douyin.RoomStreamAdaptationMessage)
	_quizAudienceStatusHandler   func(message *douyin.QuizAudienceStatusMessage)
	_hotChatHandler              func(message *douyin.HotChatMessage)
	_hotRoomHandler              func(message *douyin.HotRoomMessage)
	_audioChatHandler            func(message *douyin.AudioChatMessage)
	_roomNotifyHandler           func(message *douyin.NotifyMessage)
	_luckyBoxHandler             func(message *douyin.LuckyBoxMessage)
	_updateFanTicketHandler      func(message *douyin.UpdateFanTicketMessage)
	_screenChatHandler           func(message *douyin.ScreenChatMessage)
	_notifyEffectHandler         func(message *douyin.NotifyEffectMessage)
	_bindingGiftHandler          func(message *douyin.NotifyEffectMessage_BindingGiftMessage)
	_tempStateAreaReachHandler   func(message *douyin.TempStateAreaReachMessage)
	_growthTaskHandler           func(message *douyin.GrowthTaskMessage)
	_gameCPBaseHandler           func(message *douyin.GameCPBaseMessage)
	_unknownHandler              func(message *douyin.Message)
}

func (e *ChatHandler) chatHandler(message *douyin.ChatMessage) {
	if e._chatHandler != nil {
		e._chatHandler(message)
	}
}

func (e *ChatHandler) giftHandler(message *douyin.GiftMessage) {
	if e._giftHandler != nil {
		e._giftHandler(message)
	}
}

func (e *ChatHandler) likeHandler(message *douyin.LikeMessage) {
	if e._likeHandler != nil {
		e._likeHandler(message)
	}
}

func (e *ChatHandler) memberHandler(message *douyin.MemberMessage) {
	if e._memberHandler != nil {
		e._memberHandler(message)
	}
}

func (e *ChatHandler) socialHandler(message *douyin.SocialMessage) {
	if e._socialHandler != nil {
		e._socialHandler(message)
	}
}

func (e *ChatHandler) roomUserSeqHandler(message *douyin.RoomUserSeqMessage) {
	if e._roomUserSeqHandler != nil {
		e._roomUserSeqHandler(message)
	}
}

func (e *ChatHandler) fansclubHandler(message *douyin.FansclubMessage) {
	if e._fansclubHandler != nil {
		e._fansclubHandler(message)
	}
}

func (e *ChatHandler) controlHandler(message *douyin.ControlMessage) {
	if e._controlHandler != nil {
		e._controlHandler(message)
	}
}

func (e *ChatHandler) emojiChatHandler(message *douyin.EmojiChatMessage) {
	if e._emojiChatHandler != nil {
		e._emojiChatHandler(message)
	}
}

func (e *ChatHandler) roomStatsHandler(message *douyin.RoomStatsMessage) {
	if e._roomStatsHandler != nil {
		e._roomStatsHandler(message)
	}
}

func (e *ChatHandler) roomHandler(message *douyin.RoomMessage) {
	if e._roomHandler != nil {
		e._roomHandler(message)
	}
}

func (e *ChatHandler) ranklistHourEntranceHandler(message *douyin.RanklistHourEntranceMessage) {
	if e._ranklistHourEntranceHandler != nil {
		e._ranklistHourEntranceHandler(message)
	}
}

func (e *ChatHandler) roomRankHandler(message *douyin.RoomRankMessage) {
	if e._roomRankHandler != nil {
		e._roomRankHandler(message)
	}
}

func (e *ChatHandler) inRoomBannerHandler(message *douyin.InRoomBannerMessage) {
	if e._inRoomBannerHandler != nil {
		e._inRoomBannerHandler(message)
	}
}

func (e *ChatHandler) roomDataSyncHandler(message *douyin.RoomDataSyncMessage) {
	if e._roomDataSyncHandler != nil {
		e._roomDataSyncHandler(message)
	}
}

func (e *ChatHandler) luckyBoxTempStatusHandler(message *douyin.LuckyBoxTempStatusMessage) {
	if e._luckyBoxTempStatusHandler != nil {
		e._luckyBoxTempStatusHandler(message)
	}
}

func (e *ChatHandler) decorationModifyHandler(message *douyin.DecorationModifyMessage) {
	if e._decorationModifyHandler != nil {
		e._decorationModifyHandler(message)
	}
}

func (e *ChatHandler) linkMicAudienceKtvHandler(message *douyin.LinkMicAudienceKtvMessage) {
	if e._linkMicAudienceKtvHandler != nil {
		e._linkMicAudienceKtvHandler(message)
	}
}

func (e *ChatHandler) roomStreamAdaptationHandler(message *douyin.RoomStreamAdaptationMessage) {
	if e._roomStreamAdaptationHandler != nil {
		e._roomStreamAdaptationHandler(message)
	}
}

func (e *ChatHandler) quizAudienceStatusHandler(message *douyin.QuizAudienceStatusMessage) {
	if e._quizAudienceStatusHandler != nil {
		e._quizAudienceStatusHandler(message)
	}
}

func (e *ChatHandler) hotChatHandler(message *douyin.HotChatMessage) {
	if e._hotChatHandler != nil {
		e._hotChatHandler(message)
	}
}

func (e *ChatHandler) hotRoomHandler(message *douyin.HotRoomMessage) {
	if e._hotRoomHandler != nil {
		e._hotRoomHandler(message)
	}
}

func (e *ChatHandler) audioChatHandler(message *douyin.AudioChatMessage) {
	if e._audioChatHandler != nil {
		e._audioChatHandler(message)
	}
}

func (e *ChatHandler) roomNotifyHandler(message *douyin.NotifyMessage) {
	if e._roomNotifyHandler != nil {
		e._roomNotifyHandler(message)
	}
}

func (e *ChatHandler) luckyBoxHandler(message *douyin.LuckyBoxMessage) {
	if e._luckyBoxHandler != nil {
		e._luckyBoxHandler(message)
	}
}

func (e *ChatHandler) updateFanTicketHandler(message *douyin.UpdateFanTicketMessage) {
	if e._updateFanTicketHandler != nil {
		e._updateFanTicketHandler(message)
	}
}

func (e *ChatHandler) screenChatHandler(message *douyin.ScreenChatMessage) {
	if e._screenChatHandler != nil {
		e._screenChatHandler(message)
	}
}

func (e *ChatHandler) notifyEffectHandler(message *douyin.NotifyEffectMessage) {
	if e._notifyEffectHandler != nil {
		e._notifyEffectHandler(message)
	}
}

func (e *ChatHandler) bindingGiftHandler(message *douyin.NotifyEffectMessage_BindingGiftMessage) {
	if e._bindingGiftHandler != nil {
		e._bindingGiftHandler(message)
	}
}

func (e *ChatHandler) tempStateAreaReachHandler(message *douyin.TempStateAreaReachMessage) {
	if e._tempStateAreaReachHandler != nil {
		e._tempStateAreaReachHandler(message)
	}
}

func (e *ChatHandler) growthTaskHandler(message *douyin.GrowthTaskMessage) {
	if e._growthTaskHandler != nil {
		e._growthTaskHandler(message)
	}
}

func (e *ChatHandler) gameCPBaseHandler(message *douyin.GameCPBaseMessage) {
	if e._gameCPBaseHandler != nil {
		e._gameCPBaseHandler(message)
	}
}

func (e *ChatHandler) unknownHandler(message *douyin.Message) {
	if e._unknownHandler != nil {
		e._unknownHandler(message)
	}
}

func WithChat(fn func(*douyin.ChatMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._chatHandler = fn
	}
}

func WithGift(fn func(*douyin.GiftMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._giftHandler = fn
	}
}

func WithLike(fn func(*douyin.LikeMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._likeHandler = fn
	}
}

func WithMember(fn func(*douyin.MemberMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._memberHandler = fn
	}
}

func WithSocial(fn func(*douyin.SocialMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._socialHandler = fn
	}
}

func WithRoomUserSeq(fn func(*douyin.RoomUserSeqMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomUserSeqHandler = fn
	}
}

func WithFansclub(fn func(*douyin.FansclubMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._fansclubHandler = fn
	}
}

func WithControl(fn func(*douyin.ControlMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._controlHandler = fn
	}
}

func WithEmojiChat(fn func(*douyin.EmojiChatMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._emojiChatHandler = fn
	}
}

func WithRoomStats(fn func(*douyin.RoomStatsMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomStatsHandler = fn
	}
}

func WithRoom(fn func(*douyin.RoomMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomHandler = fn
	}
}

func WithRanklistHourEntrance(fn func(*douyin.RanklistHourEntranceMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._ranklistHourEntranceHandler = fn
	}
}

func WithRoomRank(fn func(*douyin.RoomRankMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomRankHandler = fn
	}
}

func WithInRoomBanner(fn func(*douyin.InRoomBannerMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._inRoomBannerHandler = fn
	}
}

func WithRoomDataSync(fn func(*douyin.RoomDataSyncMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomDataSyncHandler = fn
	}
}

func WithLuckyBoxTempStatus(fn func(*douyin.LuckyBoxTempStatusMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._luckyBoxTempStatusHandler = fn
	}
}

func WithDecorationModify(fn func(*douyin.DecorationModifyMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._decorationModifyHandler = fn
	}
}

func WithLinkMicAudienceKtv(fn func(*douyin.LinkMicAudienceKtvMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._linkMicAudienceKtvHandler = fn
	}
}

func WithRoomStreamAdaptation(fn func(*douyin.RoomStreamAdaptationMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomStreamAdaptationHandler = fn
	}
}

func WithQuizAudienceStatus(fn func(*douyin.QuizAudienceStatusMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._quizAudienceStatusHandler = fn
	}
}

func WithHotChat(fn func(*douyin.HotChatMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._hotChatHandler = fn
	}
}

func WithHotRoom(fn func(*douyin.HotRoomMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._hotRoomHandler = fn
	}
}

func WithAudioChat(fn func(*douyin.AudioChatMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._audioChatHandler = fn
	}
}

func WithRoomNotify(fn func(*douyin.NotifyMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._roomNotifyHandler = fn
	}
}

func WithLuckyBox(fn func(*douyin.LuckyBoxMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._luckyBoxHandler = fn
	}
}

func WithUpdateFanTicket(fn func(*douyin.UpdateFanTicketMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._updateFanTicketHandler = fn
	}
}

func WithScreenChat(fn func(*douyin.ScreenChatMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._screenChatHandler = fn
	}
}

func WithNotifyEffect(fn func(*douyin.NotifyEffectMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._notifyEffectHandler = fn
	}
}

func WithBindingGift(fn func(*douyin.NotifyEffectMessage_BindingGiftMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._bindingGiftHandler = fn
	}
}

func WithTempStateAreaReach(fn func(*douyin.TempStateAreaReachMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._tempStateAreaReachHandler = fn
	}
}

func WithGrowthTask(fn func(*douyin.GrowthTaskMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._growthTaskHandler = fn
	}
}

func WithGameCPBase(fn func(*douyin.GameCPBaseMessage)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._gameCPBaseHandler = fn
	}
}

func WithUnknown(fn func(*douyin.Message)) EventHandlerOption {
	return func(handler *ChatHandler) {
		handler._unknownHandler = fn
	}
}

func NewEventHandler(options ...EventHandlerOption) *ChatHandler {
	eventHandler := &ChatHandler{}
	for _, option := range options {
		option(eventHandler)
	}
	return eventHandler
}
