// models_extended.go
package database

type FileData struct {
	ID         int64   `json:"id"`
	Hash       string  `json:"hash"`
	FileName   string  `json:"FileName"`
	FileSize   string  `json:"FileSize"`
	SizeInGB   float64 `json:"sizeInGB"`
	DateListed string  `json:"DateListed"`
	Type       string  `json:"type"`
	Downloads  int     `json:"downloads"`
	Price      float64 `json:"price"`
}
