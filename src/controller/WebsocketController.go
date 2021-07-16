package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alpha-supsys/go-common/app/rest"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的CORS 跨域请求，正式环境可以关闭
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketController struct {
	rest.Controller
	m map[string]*websocket.Conn
}

func NewWebsocketController() *WebsocketController {
	return &WebsocketController{
		m: make(map[string]*websocket.Conn),
	}
}

func (s *WebsocketController) GetRoute() rest.RouteMap {
	items := []*rest.RouteItem{
		{Path: "/api/beta/ws", HandleFunc: s.ws_connect, Method: "GET"},
	}

	return rest.NewRouteMap(items...)
}

func (s *WebsocketController) ws_connect(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级为websocket失败", err.Error())
		return
	}
	fmt.Println("ssl-client-cert:", r.Header.Get("ssl-client-cert"))
	// s.m[fmt.Sprintf("%p", wsConn)] =
	go s.syncState(wsConn)
	err = s.msgHander(wsConn)
	delete(s.m, fmt.Sprintf("%p", wsConn))
	if err != nil {
		log.Println("msg:", err.Error())
		return
	}
}

func (s *WebsocketController) syncState(wsConn *websocket.Conn) {
	for {
		// 修改以下内容把客户端传递的消息传递给处理程序
		err := wsConn.WriteJSON(s.m)
		if err != nil {
			log.Println("发送消息给客户端出现错误", err.Error())
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *WebsocketController) msgHander(wsConn *websocket.Conn) error {
	for {
		msgType, data, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("获取消息出现错误", err.Error())
			return err
		}
		log.Println("接收到消息", string(data))
		// 修改以下内容把客户端传递的消息传递给处理程序
		err = wsConn.WriteMessage(msgType, data)
		if err != nil {
			log.Println("发送消息给客户端出现错误", err.Error())
			return err
		}
	}
}
