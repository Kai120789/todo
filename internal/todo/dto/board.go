package dto

type GetBoardDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PostBoardDto struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}
