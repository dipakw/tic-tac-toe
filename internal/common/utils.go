package common

import (
	"encoding/json"
	"math/rand"
)

func GenerateAlphaNumId(size int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	id := make([]byte, size)

	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}

	return string(id)
}

func DecodeAnyAs[T any](input any) (*T, error) {
	data, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	var v T

	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
