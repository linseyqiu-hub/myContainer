package IDGeneration

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID(numBytes int) (string, error) {
	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
