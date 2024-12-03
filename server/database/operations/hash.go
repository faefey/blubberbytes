package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Takes a file at located filePath returns the sha256 encoding of its contents in hex form
// To test it out, try hashing cmd.txt and then finding the hash checksum here:
// https://emn178.github.io/online-tools/sha256_checksum.html
func HashFile(filePath string) (string, error) {
	fileContent, err := os.Open(filePath)
	if err != nil {
		return "Error hashing", err
	}
	defer fileContent.Close()

	h := sha256.New()
	if _, err := io.Copy(h, fileContent); err != nil {
		return "Error hashing", err
	}

	hash := hex.EncodeToString(h.Sum(nil))
	fmt.Println("Hash of file at " + filePath + ": " + hash)
	return hash, nil
}
