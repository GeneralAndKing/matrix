package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kernel/internal/message"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func _ws(c *gin.Context, namespace string) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update ws: %w", err))
		return
	}
	go func() {
		for {
			_, _, err = ws.ReadMessage()
			if err != nil {
				_ = ws.Close()
				return
			}
		}

	}()
	err = message.Subscribe(c, namespace, func(m message.Message) error {
		return ws.WriteJSON(m)
	})
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("ws broke the link: %w", err))
		return
	}

}

func Message(c *gin.Context) {
	_ws(c, message.WS)

}

func Business(c *gin.Context) {
	_ws(c, message.BUSINESS)
}
