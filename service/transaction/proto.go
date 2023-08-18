package transaction

type ChargeData struct {
	Address string `json:"address"`
	TxHash  string `json:"txHash"`
	Amount  int    `json:"amount"`
}
