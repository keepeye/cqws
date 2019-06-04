package cqws

import (
	"bytes"
	"cqwss/cqws/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type PostType int

const (
	PostTypeMessage PostType = iota
	PostTypeNotice
	PostTypeRequest
	PostTypeUnknown
)

// 判断请求类型
func detectPostType(body []byte) PostType {
	if bytes.Contains(body, []byte("\"post_type\":\"message\"")) {
		return PostTypeMessage
	}
	if bytes.Contains(body, []byte("\"post_type\":\"notice\"")) {
		return PostTypeNotice
	}
	if bytes.Contains(body, []byte("\"post_type\":\"request\"")) {
		return PostTypeRequest
	}
	return PostTypeUnknown
}

func parseMessage(body []byte) *models.Message {
	var m models.Message
	err := json.Unmarshal(body, &m)
	if err != nil {
		logrus.WithField("body", string(body)).WithField("error", err.Error()).Error("消息解构失败")
		return nil
	}
	return &m
}
