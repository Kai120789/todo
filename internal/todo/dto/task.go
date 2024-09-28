package dto

type GetTaskDto struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type PostTaskDto struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Board_id    string `json:"board_id"`
	Status_id   string `json:"status_id"`
	User_id     string `json:"user_id"`
}
