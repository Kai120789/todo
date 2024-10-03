package dto

type GetTokenDto struct {
	RefreshToken string `json:"refresh_token"`
}

type PostTokenDto struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}
