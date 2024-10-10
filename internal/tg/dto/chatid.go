package dto

type ChatID struct {
	Username string `json:"tg_name"`
	ChatID   int64  `json:"chat_id"`
}
