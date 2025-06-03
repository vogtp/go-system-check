package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Calc calculate the file hash
func Calc() (string, error) {
	//defer func(t time.Time) {fmt.Printf("Hashing took %v\n", time.Since(t))}(time.Now())
	f, err := os.Open(os.Args[0])
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
