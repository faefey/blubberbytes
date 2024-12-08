package models

// Struct (not a table) for Wallet
type Wallet struct {
	Address        string  `json:"address"`
	CurrentBalance float64 `json:"currentBalance"`
	PendingBalance float64 `json:"pendingBalance"`
}
