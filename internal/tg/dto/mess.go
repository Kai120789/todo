package dto

type MessDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusId    uint   `json:"status_id"`
	ChatId      int64
}
