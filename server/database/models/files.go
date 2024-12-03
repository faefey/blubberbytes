package models

// Table for Storing
type Storing struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Date      string `json:"date"`
}

// Table for Hosting
type Hosting struct {
	Hash  string  `json:"hash"`
	Price float64 `json:"price"`
}

// Struct (not a table) for Hosting joined with Storing
type JoinedHosting struct {
	Hash      string  `json:"hash"`
	Name      string  `json:"name"`
	Extension string  `json:"extension"`
	Size      int64   `json:"size"`
	Path      string  `json:"path"`
	Date      string  `json:"date"`
	Price     float64 `json:"price"`
}

// Table for Sharing
type Sharing struct {
	Hash     string `json:"hash"`
	Password string `json:"password"`
}

// Struct (not a table) for Sharing joined with Storing
type JoinedSharing struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Date      string `json:"date"`
	Password  string `json:"password"`
}

// Table for saved files
type Saved struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}
