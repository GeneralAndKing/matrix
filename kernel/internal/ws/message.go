package ws

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"go.uber.org/zap"
	"net/http"
)

type MessageType uint

var (
	server *socketio.Server
)

const (
	MessageDEBUG MessageType = iota
	MessageINFO
	MessageWARN
	MessageERROR
	messageNameSpace = "/message"
)

func Init() {
	server = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{CheckOrigin: func(r *http.Request) bool {
				return true
			}},
			&websocket.Transport{CheckOrigin: func(r *http.Request) bool {
				return true
			}}}})
	//Add all connected user to a room, in example? "bcast"
	server.OnConnect(messageNameSpace, func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})
	server.OnDisconnect(messageNameSpace, func(s socketio.Conn, s2 string) {

	})
	go func() {
		if err := server.Serve(); err != nil {
			zap.L().Fatal("socketio listen error", zap.Error(err))
		}
	}()
}
func Open() error {
	return server.Serve()
}
func Close() error {
	return server.Close()
}

func Handler() *socketio.Server {
	return server
}

func BroadcastToMessage(_type MessageType, s string) {
	marshal, _ := json.Marshal(map[string]interface {
	}{"status": _type, "message": s})
	server.BroadcastToNamespace(messageNameSpace, "receive", marshal)
}
