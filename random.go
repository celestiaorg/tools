// This package contains the tools that frequently used by other packages
package tools

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// This function receives a pair of min and max float64 numbers
// and generates a random number in that range (inclusive)
func RandomNumberF(rangeLower float64, rangeUpper float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rangeLower + rand.Float64()*(rangeUpper-rangeLower)
}

// This function receives a pair of min and max int64 numbers
// and generates a random number in that range (inclusive)
func RandomNumberI(rangeLower int64, rangeUpper int64) int64 {

	rand.Seed(time.Now().UnixNano())
	return rangeLower + rand.Int63n(rangeUpper-rangeLower+1)
}

// It generates a random boolean status
// the `chancePercent` should be between 0 and 100
func RandomBool(chancePercent int) bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100+1) <= chancePercent
}

// This function receives the length of a string
// and generates a random string
func RandomString(length int) string {

	rStr := ""
	for len(rStr) < length {
		rand.Seed(time.Now().UnixNano())

		rn := rand.Int63()
		rStr += fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d", rn))))

	}
	output := strings.Trim(base64.StdEncoding.EncodeToString([]byte(rStr)), "=")
	return output[:length]
}
