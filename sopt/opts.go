package sopt

import "cqwss/cqws"

func WsPath(path string) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.WsPath = path
	}
}

func Host(host string) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.Host = host
	}
}

func Port(port string) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.Port = port
	}
}

func MessageHandler(f cqws.MessageHandleFunc) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.MessageHandler = f
	}
}

func SetAccessToken(token string) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.AccessToken = token
	}
}

func SetWhiteList(qqList []string) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.WhiteList = qqList
	}
}

func SetBlackList(qqList []string) cqws.ServerOpt {
	return func(s *cqws.Server) {
		s.BlackList = qqList
	}
}
