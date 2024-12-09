package models

// Table for WalletInfo
type WalletInfo struct {
	Address        string `json:"address"`
	PubPassphrase  string `json:"pubPassphrase"`
	PrivPassphrase string `json:"privPassphrase"`
}

// Struct (not a table) for Wallet
type Wallet struct {
	Address        string  `json:"address"`
	CurrentBalance float64 `json:"currentBalance"`
	PendingBalance float64 `json:"pendingBalance"`
}
