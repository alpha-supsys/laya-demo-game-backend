package controller

import (
	"crypto/sha1"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

func cert2id(str string) string {
	decode_str, _ := url.QueryUnescape(str)
	p, _ := pem.Decode([]byte(decode_str))
	cry := sha1.New()
	cry.Write(p.Bytes)

	return fmt.Sprintf("%x", cry.Sum(nil))
}

func (s *WebsocketController) ws_connect(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级为websocket失败", err.Error())
		return
	}
	header_cert_str := r.Header.Get("ssl-client-cert")
	cert_id := cert2id(header_cert_str)
	fmt.Println("ssl-client-cert:", header_cert_str)
	s.m[cert_id] = wsConn
	go s.syncState(wsConn)
	err = s.msgHander(wsConn)
	if err != nil {
		log.Println("msg:", err.Error())
		return
	}
	if wsConn, ok := s.m[cert_id]; ok {
		wsConn.Close()
		delete(s.m, cert_id)
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
