package models

// Table for upload history
type Uploads struct {
	Id        int64  `json:"id"`
	Date      string `json:"date"`
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

// Table for download history
type Downloads struct {
	Id        int64   `json:"id"`
	Date      string  `json:"date"`
	Hash      string  `json:"hash"`
	Name      string  `json:"name"`
	Extension string  `json:"extension"`
	Size      int64   `json:"size"`
	Price     float64 `json:"price"`
}

// Table for transaction history
type Transactions struct {
	Id      int64   `json:"id"`
	Date    string  `json:"date"`
	Wallet  string  `json:"wallet"`
	Amount  float64 `json:"amount"`
	Balance float64 `json:"balance"`
}

// Put this on hold until proxies are implemented
// Table for proxy history
type Proxies struct {
}
