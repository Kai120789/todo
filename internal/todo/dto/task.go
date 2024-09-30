package dto

type GetTaskDto struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type PostTaskDto struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BoardId     string `json:"board_id"`
	StatusId    string `json:"status_id"`
	UserId      string `json:"user_id"`
}
