package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// Claims 自定义claims
type Claims struct {
	Phone string
	jwt.StandardClaims
}

var key = []byte("abcdefg.654321")

// CreateToken 生成jwt-token
func CreateToken(phone string) (string, error) {

	expiresTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Phone: phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "filbox",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
		return key, nil
	})
	return token, claims, err
}
