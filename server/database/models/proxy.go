package models

// Table for Proxy
type Proxy struct {
	IP     string  `json:"ip"`
	Port   string  `json:"port"`
	Rate   float64 `json:"rate"`
	Wallet string  `json:"wallet"`
}
