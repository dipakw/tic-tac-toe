package common

import "math/rand"

func GenerateAlphaNumId(size int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	id := make([]byte, size)

	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}

	return string(id)
}
