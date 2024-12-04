package models

type Wallet struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}
