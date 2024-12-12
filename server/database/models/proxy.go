package models

// Table for Proxy
type Proxy struct {
	IP     string  `json:"ip"`
	Rate   float64 `json:"rate"`
	Node   string  `json:"node"`
	Wallet string  `json:"wallet"`
}

// Table for ProxyLogs
type ProxyLogs struct {
	Id    string `json:"id"`
	IP    string `json:"ip"`
	Bytes int64  `json:"bytes"`
	Time  int64  `json:"time"`
}

// Table for IPtoNode
type IPtoNode struct {
	IP   string `json:"ip"`
	Node string `json:"node"`
}

// Struct (not a table) for ProxyBill
type ProxyBill struct {
	IP     string  `json:"ip"`
	Rate   float64 `json:"rate"`
	Bytes  int64   `json:"bytes"`
	Amount float64 `json:"amount"`
	Wallet string  `json:"wallet"`
}
