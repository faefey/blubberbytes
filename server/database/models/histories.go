package models

// Table for Uploads
type Uploads struct {
	Id        int64  `json:"id"`
	Date      string `json:"date"`
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

// Table for Downloads
type Downloads struct {
	Id        int64   `json:"id"`
	Date      string  `json:"date"`
	Hash      string  `json:"hash"`
	Name      string  `json:"name"`
	Extension string  `json:"extension"`
	Size      int64   `json:"size"`
	Price     float64 `json:"price"`
}

// Struct (not a table) for Transactions
type Transactions struct {
	Id            string  `json:"id"`
	Date          string  `json:"date"`
	Wallet        string  `json:"wallet"`
	Amount        float64 `json:"amount"`
	Category      string  `json:"category"`
	Fee           float64 `json:"fee"`
	Confirmations int64   `json:"confirmations"`
}

// Put this on hold until proxies are implemented
// Table for Proxies
type Proxies struct {
}
