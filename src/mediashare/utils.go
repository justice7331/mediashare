package main

import (
	"crypto/rand"
	"encoding/base64"
	"math"
)

func randomString(length int) string {
	buff := make([]byte, int(math.Round(float64(length)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:length]
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
