package dto

type Dto struct {
	ChatID       int64  `json:"chat_id"`
	Message      string `json:"message"`
	MessageEnded string `json:"message_ended"`
}
