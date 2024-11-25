package models

// Table for added files
type Storing struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Date      string `json:"date"`
}

// Join on file hash with Storing table
type Hosting struct {
	Hash  string  `json:"hash"`
	Price float64 `json:"price"`
}

// Join on file hash with Storing table
type Sharing struct {
	Hash     string `json:"hash"`
	Password string `json:"password"`
}

// Table for saved files
type Saved struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}
