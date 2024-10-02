package dto

type GetBoardDto struct {
	Name string `json:"name"`
}

type PostBoardDto struct {
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}
