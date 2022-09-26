package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func ZeroFill(value interface{}, count int) string {
	return strings.Repeat("0", count-len(fmt.Sprintf("%v", value))) + fmt.Sprintf("%v", value)
}

// Sqrt10 求制定数的10次方根
func Sqrt10(n int) int {
	if n < 10 {
		return 1
	}
	return Sqrt10(n / 10)
}

// 求制定数的制定次方根
func SqrtN(n, m int) int {
	if n < m {
		return 1
	}
	return m * SqrtN(n/m, m)
}

const randomStringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var maxRandomStringCharsLength = len(randomStringChars)

// GetRandomString return a securely generated random string
func GetRandomString(length int) string {
	if length > maxRandomStringCharsLength {
		length = maxRandomStringCharsLength
	}
	var builder strings.Builder
	builder.Grow(length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		index := rand.Intn(maxRandomStringCharsLength)
		builder.WriteString(string(randomStringChars[index]))
	}
	return builder.String()
}

// HidePhone 将手机号码中间四位改成****
func HidePhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
