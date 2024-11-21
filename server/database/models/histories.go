package database

// Table for upload history
type Uploads struct {
	Id        int64  `json:"id"`
	Date      string `json:"date"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Hash      string `json:"hash"`
}

// Table for download history
type Downloads struct {
	Id        int64   `json:"id"`
	Date      string  `json:"date"`
	Name      string  `json:"name"`
	Extension string  `json:"extension"`
	Size      int64   `json:"size"`
	Price     float64 `json:"price"`
	Hash      string  `json:"hash"`
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
