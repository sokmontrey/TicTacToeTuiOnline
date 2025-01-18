package pkg

import (
	"math/rand"
	"strconv"
)

func GenerateId(numDigits int) string {
	id := ""
	for i := 0; i < numDigits; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}
	return id
}
