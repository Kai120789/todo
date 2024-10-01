package dto

type GetUserDto struct {
	ID       string `json:"id"`
	Username string `json:"title"`
}

type PostUserDto struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
