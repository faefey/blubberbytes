// models.go
package database

// FileMetadata stores both file metadata and provider-specific information
type FileMetadata struct {
	ID        int64   `json:"id"`         // Unique ID for each file entry
	FileSize  int64   `json:"file_size"`  // Size of the file
	Extension string  `json:"extension"`  // File extension
	FileName  string  `json:"file_name"`  // Name of the file provided by this peer
	FilePrice float64 `json:"file_price"` // Price set by the provider
	FileHash  string  `json:"file_hash"`  // Unique hash to identify the file
}
