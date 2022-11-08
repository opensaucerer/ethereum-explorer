package transaction

type Transaction struct {
	BlockNumber int64   `json:"block_number"`
	Hash        string  `json:"hash"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Amount      float64 `json:"amount"`
	Nonce       uint64  `json:"nonce"`
	BlockTime   int64   `json:"block_time"`
}
