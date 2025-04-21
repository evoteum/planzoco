package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	idLength = 6  // 62^6 possible IDs should be enough for up to 36m customers.
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GenerateID() (string, error) {
	return gonanoid.Generate(alphabet, idLength)
}
