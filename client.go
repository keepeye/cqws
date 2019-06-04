package cqws

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// 发送消息超时时间
	writeWait = 3 * time.Second

	// 客户端发来的心跳超时时间，超过该时间，终止连接
	pongWait = 60 * time.Second

	// 服务器主动向客户端发送心跳包的间隔时间，必须比pongWait小
	pingPeriod = 20 * time.Second

	// 消息体最大字节长度
	maxMessageSize = 4096
)

type Client struct {
	server *Server
	hub    *Hub
	conn   *websocket.Conn
	// 用于发送消息的管道
	send chan []byte

	// 客户端 qq 号
	QQ string `json:"qq"`
}

// 发送消息
func (c *Client) SendMessage(toUserID uint, text string) error {
	m := make(map[string]interface{})
	m["action"] = "send_private_msg"
	m["params"] = map[string]interface{}{
		"user_id": toUserID,
		"message": text,
	}
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	c.send <- b
	return nil
}

// 在单独的goruntine中运行，接收消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, body, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
		}
		switch detectPostType(body) {
		case PostTypeMessage:
			if c.server.MessageHandler != nil {
				message := parseMessage(body)
				c.server.MessageHandler(c, message)
			}
		}
	}
}

// 在单独的goroutine中运行，从管道中获取消息并发送给客户端
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		// 优先处理心跳
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 数据库主动关闭连接
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// 向客户端发送消息
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 1024,
}

func serveWs(server *Server, hub *Hub, w http.ResponseWriter, r *http.Request) error {
	if "Universal" != r.Header.Get("X-Client-Role") {
		return errors.New("cqhttp插件配置错误，ws_reverse_use_universal_client不为true")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.New("websocket upgrader error:" + err.Error())
	}
	client := &Client{server: server, hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.QQ = r.Header.Get("X-Self-ID")
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
	return nil
}
