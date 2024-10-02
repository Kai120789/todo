package dto

type GetUserDto struct {
	Username string `json:"username"`
}

type PostUserDto struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
