package utils

import (
    "errors"
    "math/big"
    "crypto/rand"
)

var tokenChars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateToken(length uint64) (*string, error) {
    chars := make([]rune, length)

    for i := range chars {
        charId, err := rand.Int(rand.Reader, big.NewInt(int64(len(tokenChars))))
        if err != nil {
            return nil, errors.New("Failed to generate random number")
        }
        chars[i] = tokenChars[charId.Int64()]
    }

    token := string(chars)
    return &token, nil
}
