package business

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kernel/internal/database"
	"kernel/internal/model"
	"kernel/internal/model/enum"
	"kernel/internal/ws"
	"time"
)

func Init(ctx context.Context) {
	go func() {
		//等待5秒 后台所有功能初始完毕
		time.Sleep(5 * time.Second)

	}()
}

func publishDouyinCreationHandle(ctx context.Context) {
	var douyinCreation model.DouyinCreation
	err := database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if tx := db.Where("status = ? ", enum.PendingCreationStatus).Order("created_at").First(&douyinCreation); tx.Error != nil {
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to found pending douyinCreation")
		}
		douyinCreation.Status = enum.RunningCreationStatus
		if tx := db.Save(&douyinCreation); tx.Error != nil {
			return fmt.Errorf("failed to save running douyin creation")
		}
		return nil
	})
	if err != nil {
		zap.L().Warn("failed to publish douyinCreation", zap.Error(err))
	}
	ws.BroadcastToMessage(ws.MessageINFO, fmt.Sprintf("[%d] 正在发布作品 %s", douyinCreation.DouyinUserId, douyinCreation.Title))
}
