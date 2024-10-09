package dto

type GetTaskDto struct {
	Title string `json:"title"`
}

type PostTaskDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	BoardId     string `json:"board_id"`
	StatusId    uint   `json:"status_id"`
	UserId      string `json:"user_id"`
}
