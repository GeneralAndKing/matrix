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

func publishDouyinWorkHandle(ctx context.Context) {
	var douyinWork model.DouyinWork
	err := database.Sqlite3Transaction(ctx, func(db *gorm.DB) error {
		if tx := db.Where("status = ? ", enum.PendingWorkStatus).Order("created_at").First(&douyinWork); tx.Error != nil {
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to found pending douyinWork")
		}
		douyinWork.Status = enum.RunningWorkStatus
		if tx := db.Save(&douyinWork); tx.Error != nil {
			return fmt.Errorf("failed to save running douyinWork")
		}
		return nil
	})
	if err != nil {
		zap.L().Warn("failed to publish douyinWork", zap.Error(err))
	}
	ws.BroadcastToMessage(ws.MessageINFO, fmt.Sprintf("[%d] 正在发布作品 %s", douyinWork.DouyinUserId, douyinWork.Title))
}
