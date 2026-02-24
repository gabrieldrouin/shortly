package shortcode

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length   = 7
)

var alphabetLen = big.NewInt(int64(len(alphabet)))

func Generate() (string, error) {
	code := make([]byte, length)
	for i := range code {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		code[i] = alphabet[idx.Int64()]
	}
	return string(code), nil
}
