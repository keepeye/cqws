package models

type Message struct {
	BaseInfo
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageID   uint   `json:"message_id"`
	UserID      uint   `json:"user_id"`
	RawMessage  string `json:"raw_message"`
	Font        uint   `json:"font"`
	Sender      struct {
		UserID   uint   `json:"user_id"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      uint   `json:"age"`
	} `json:"sender"`
	Message []MessageSegment `json:"message"`
}

type MessageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}
