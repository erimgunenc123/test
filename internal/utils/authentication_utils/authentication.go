package authentication_utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"genericAPI/api/api_config"
	"genericAPI/api/database_connection"
	"genericAPI/internal/customErrors"
	"genericAPI/internal/models/refresh_token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"log"
	mathrand "math/rand"
	"strconv"
	"strings"
	"time"
)

const hashIterationCount = 10000
const accessTokenDurationSeconds = 3600    // 1 hour
var applicationIdentifier = mathrand.Int() // just to make sure uuid's won't collide across all instances

func CreateAccessToken(userId uint64, publicId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"public_id": publicId,
		"user_id":   userId,
		"exp":       time.Now().Add(time.Second * accessTokenDurationSeconds).UnixMilli(),
	})

	tokenString, err := token.SignedString([]byte(api_config.Config.App.Secret))

	return tokenString, err
}

func CreatePublicID() string {
	newUuid := uuid.New()
	return fmt.Sprintf("%s%d", newUuid.String(), applicationIdentifier)

}

// ValidateAccessToken validates and returns the user id
func ValidateAccessToken(tokenString string) (userid *uint64, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, customErrors.ErrUnexpectedSigningMethod
		}
		return []byte(api_config.Config.App.Secret), nil
	})
	if err != nil {
		return nil, customErrors.ErrInvalidAccessToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var publicId, userId string
		var expiration int64
		if publicId_, ok := claims["public_id"]; ok {
			publicId = publicId_.(string)
		}
		if userId_, ok := claims["user_id"]; ok {
			userId = userId_.(string)
		}
		if expiration__, ok := claims["exp"]; ok {
			if expiration_, ok := expiration__.(string); ok {
				expiration, _ = strconv.ParseInt(expiration_, 10, 64)
			}
		}
		if publicId == "" || userId == "" || expiration == 0 {
			return nil, customErrors.ErrInvalidAccessToken
		}
		if time.Now().After(time.UnixMilli(expiration)) {
			return nil, customErrors.ErrAccessTokenExpired
		}
		uintUserId, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			return nil, customErrors.ErrInvalidAccessToken
		}
		return &uintUserId, nil

	}
	return nil, customErrors.ErrInvalidAccessToken
}

func CreateRefreshToken(userId uint64) string {
	concatStr := []byte(strconv.FormatUint(userId, 10) + strconv.FormatInt(time.Now().UnixMilli(), 10))
	return base64.StdEncoding.EncodeToString(concatStr)
}

func ValidateRefreshToken(refreshToken string, userId uint64) bool {
	var token refresh_token.RefreshToken
	err := database_connection.DB.Table("refresh_token").Where("user_id = ?", userId).First(&token).Error
	if err != nil {
		return false
	}
	return token.Token == refreshToken
}

func HashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	return _hashString(password, saltBase64, hashIterationCount, 32)
}

func _hashString(targetStr string, salt string, iterationCount int, keyLen int) string {
	hashedPw := pbkdf2.Key([]byte(targetStr), []byte(salt), iterationCount, keyLen, sha256.New)
	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", hashIterationCount, salt, base64.StdEncoding.EncodeToString(hashedPw))
}

func ComparePassword(hashedPw string, rawPw string) bool {
	splitted := strings.Split(hashedPw, "$")
	if len(splitted) != 4 {
		log.Printf("Invalid stored password: %s", hashedPw)
		return false
	}
	salt := splitted[2]
	iterationCount, err := strconv.Atoi(splitted[1])
	if err != nil {
		log.Printf("Invalid iteration count in password: %s", hashedPw)
		return false
	}
	hashedRawPw := _hashString(rawPw, salt, iterationCount, 32)
	return hashedRawPw == hashedPw
}
