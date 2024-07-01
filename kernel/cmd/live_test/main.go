package main

import (
	"kernel/pkg/douyin_live"
	"kernel/pkg/douyin_live/generated/douyin"
)

func main() {
	monitor, err := douyin_live.New("746786991471")
	if err != nil {
		panic(err)
	}
	monitor.Start(douyin_live.NewEventHandler(douyin_live.WithChat(func(message *douyin.ChatMessage) {
		println(message.Content)
	}), douyin_live.WithRoomUserSeq(func(message *douyin.RoomUserSeqMessage) {
		println("userSeq人数", message.Total)
	}), douyin_live.WithRoomStats(func(message *douyin.RoomStatsMessage) {
		println("roomStats人数", message.Total)

	})))
}
