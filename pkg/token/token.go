package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExp       = time.Hour * 3
	TokenSecretKey = "supersecretkey"
)

type Claims struct {
	jwt.RegisteredClaims
}

type Token interface {
	Create() (string, error)
	Check(tokenString string) bool
}

func Create() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp))},
		},
	)
	signedString, err := token.SignedString([]byte(TokenSecretKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func Check(tokenString string) (bool, error) {
	data := &Claims{}

	if _, err := jwt.ParseWithClaims(
		tokenString, data,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(TokenSecretKey), nil
		},
	); err != nil {
		return false, err
	}

	return true, nil
}
