package dto

type GetBoardDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PostBoardDto struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	User_id string `json:"user_id"`
}
