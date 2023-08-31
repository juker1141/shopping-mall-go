package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomID generates a random ID
func RandomID() int32 {
	return int32(RandomInt(1, 100))
}

// RandomAccount generates a random account
func RandomAccount() string {
	return RandomString(8)
}

// RandomName generates a random name
func RandomName() string {
	return RandomString(6)
}

// RandomPermission generates a random permission
func RandomPermission() string {
	return RandomString(8)
}

// RandomRole generates a random role
func RandomRole() string {
	return RandomString(4)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(8))
}

func RandomGender() int32 {
	return int32(RandomInt(1, 3))
}

func RandomPhone() string {
	return RandomString(10)
}

func RandomAddress() string {
	return RandomString(20)
}

func RandomPostCode() string {
	return strconv.Itoa(int(RandomInt(101, 999)))
}

func RandomPrice() int32 {
	return int32(RandomInt(1, 1000))
}
