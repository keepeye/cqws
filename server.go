package cqws

import (
	"cqwss/cqws/models"
	"fmt"
	"log"
	"net/http"
)

type ServerOpt func(s *Server)
type MessageHandleFunc func(client *Client, message *models.Message)

type Server struct {
	WsPath string
	Host   string
	Port   string
	// 消息处理函数
	MessageHandler MessageHandleFunc
}

func NewServer(opts ...ServerOpt) *Server {
	s := &Server{
		WsPath: "/ws",
		Host:   "0.0.0.0",
		Port:   "2019",
	}
	for _, f := range opts {
		f(s)
	}
	return s
}

func (s *Server) Listen() {
	hub := newHub()
	go hub.run()
	http.HandleFunc(s.WsPath, func(w http.ResponseWriter, r *http.Request) {
		err := serveWs(s, hub, w, r)
		if err != nil {
			log.Println(err)
		}
	})
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	log.Println("Listen on ", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
