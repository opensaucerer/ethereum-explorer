package transaction

// CREATE TABLE transactions(
//     block_number INT PRIMARY KEY,
//     hash CHAR(32) NOT NULL,
//     from CHAR(32) NOT NULL,
//     to CHAR(32) NOT NULL,
//     amount NUMERIC NOT NULL,
//     nonce INT NOT NULL
// );

type Transaction struct {
	BlockNumber int64   `json:"block_number"`
	Hash        string  `json:"hash"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Amount      float64 `json:"amount"`
	Nonce       uint64  `json:"nonce"`
	BlockTime   uint64  `json:"block_time"`
}
