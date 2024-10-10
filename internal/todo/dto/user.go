package dto

type GetUserDto struct {
	Username string `json:"username"`
}

type PostUserDto struct {
	Username     string `json:"username"`
	TgName       string `json:"tg_name"`
	ChatID       int64  `json:"chat_id"`
	PasswordHash string `json:"password"`
}
