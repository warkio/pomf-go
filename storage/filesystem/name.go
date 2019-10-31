package filesystem

import (
	"crypto/rand"
	"encoding/base32"
)

var nameEncoding = base32.NewEncoding("bcdfghjkmnpqrstwxyzBCDFG23456789").WithPadding('Z')

// RandomName generates a random string suitable as file name, of length `n`.
func RandomName(n int) (string, error) {
	buf := make([]byte, n)

	var err error

	_, err = rand.Read(buf)
	if err != nil {
		return "", err
	}

	name := nameEncoding.EncodeToString(buf)[:n]

	return name, nil
}
