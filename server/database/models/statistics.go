package models

// Struct (not a table) for Statistics
type Statistics struct {
	StoringNum  int64 `json:"storingNum"`
	StoringSize int64 `json:"storingSize"`
	HostingNum  int64 `json:"hostingNum"`
	HostingSize int64 `json:"hostingSize"`
	SharingNum  int64 `json:"sharingNum"`
	SharingSize int64 `json:"sharingSize"`
	SavedNum    int64 `json:"savedNum"`
	SavedSize   int64 `json:"savedSize"`
}
