package divination

type CreateDivination struct {
	UserAddress string `json:"user_address" binding:"required"`
	Card1       int    `json:"card_1"`
	Card1Status bool   `json:"card_1_status"`
	Card2       int    `json:"card_2"`
	Card2Status bool   `json:"card_2_status"`
	Card3       int    `json:"card_3"`
	Card3Status bool   `json:"card_3_status"`
	Content     string `json:"content"`
}
