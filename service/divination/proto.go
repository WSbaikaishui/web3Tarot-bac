package divination

type CreateDivination struct {
	UserID      int    `json:"user_id" binding:"required"`
	Card1       int    `json:"card_1"`
	Card1Status bool   `json:"card_1_status"`
	Question    string `json:"quesiton"`
	Content     string `json:"content"`
}
