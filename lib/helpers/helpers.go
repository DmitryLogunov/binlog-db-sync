package helpers

import (
	"math/rand"
	"os"
	"time"
)

func stringWithCharset(length int, charset string) string {	
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
  return stringWithCharset(length, charset)
}


func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
    return value
	}

	return defaultVal
}