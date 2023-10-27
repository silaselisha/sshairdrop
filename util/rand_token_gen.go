package util

import (
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyz@#&*!?~"

func RandTokenGenerator(size int) string {
	length := len(letters)

	var builder strings.Builder

	for i := 0; i < size; i++ {
		char := letters[rand.Int31n(int32(length))]

		err := builder.WriteByte(char)
		if err != nil {
			return ""
		}
	}
	return builder.String()
}