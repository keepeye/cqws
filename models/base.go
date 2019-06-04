package models

import "time"

// 把BaseInfo嵌入到每个消息结构体中

type BaseInfo struct {
	PostType string `json:"post_type"`
	Time     int    `json:"time"`
	SelfID   int    `json:"self_id"`
}

func (b *BaseInfo) GetTimeString() string {
	return time.Unix(int64(b.Time), 0).Format("2006-01-02 15:04:05")
}
