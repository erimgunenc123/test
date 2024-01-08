package authentication_utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"strconv"
	"strings"
)

const hashIterationCount = 10000

func HashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	return _hashPassword(password, saltBase64, hashIterationCount, 32)
}

func _hashPassword(password string, salt string, iterationCount int, keyLen int) string {
	hashedPw := pbkdf2.Key([]byte(password), []byte(salt), iterationCount, keyLen, sha256.New)
	return fmt.Sprintf("pbkdf2_sha256$%s$%d$%s", salt, hashIterationCount, hashedPw)
}

func ComparePassword(hashedPw string, rawPw string) bool {
	splitted := strings.Split(hashedPw, "$")
	if len(splitted) != 4 {
		log.Printf("Invalid stored password: %s", hashedPw)
		return false
	}
	salt := splitted[1]
	iterationCount, err := strconv.Atoi(splitted[2])
	if err != nil {
		log.Printf("Invalid iteration count in password: %s", hashedPw)
		return false
	}
	keyLen := len(splitted[3])
	hashedRawPw := _hashPassword(rawPw, salt, iterationCount, keyLen)
	return hashedRawPw == hashedPw
}
