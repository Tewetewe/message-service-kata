package utils

import (
	"math"
	"math/rand"
	"time"
)

// RandNumericString function random string with numeric material
func RandNumericString(length int) string {
	material := "0123456789"
	e, r := len(material), ""

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		c := int(math.Floor(rand.Float64() * float64(e)))
		r += string([]rune(material)[c])
	}
	return r
}
