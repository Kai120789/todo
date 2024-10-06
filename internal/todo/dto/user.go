package dto

type GetUserDto struct {
	Username string `json:"username"`
}

type PostUserDto struct {
	Username     string `json:"username"`
	TgName       string `json:"tg_name"`
	PasswordHash string `json:"password_hash"`
}
