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
	WsPath      string
	Host        string
	Port        string
	AccessToken string
	// 允许连接的客户端QQ白名单，默认为nil，允许所有
	WhiteList []string
	// 白名单为nil，黑名单生效
	BlackList []string
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
		qq := r.Header.Get("X-Self-ID")
		// 校验access_token
		if s.AccessToken != "" && string([]byte(r.Header.Get("Authorization"))[6:]) != s.AccessToken {
			log.Println("access_token_error:qq=", qq)
			return
		}
		// 客户端qq号白名单黑名单检测
		var allow = true
		if s.WhiteList != nil {
			allow = false
			for _, v := range s.WhiteList {
				if v == qq {
					allow = true
					break
				}
			}
			if !allow {
				log.Println("connection denied from whitelist:qq=", qq)
				return
			}
		} else {
			if s.BlackList != nil {
				for _, v := range s.BlackList {
					if v == qq {
						allow = false
						break
					}
				}
			}
		}
		// 建立websocket连接
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
