package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashList(list []string) string {
	builder := ""
	for _, s := range list {
		builder += s
	}
	return Hash(builder)
}

func Hash(stringToHash string) string {
	h := sha1.New()
	h.Write([]byte(stringToHash))
	bs := h.Sum(nil)
	sh := fmt.Sprintf("%x", bs)
	return sh
}
