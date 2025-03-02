package token

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{
		secretKey: secretKey,
	}
}

func (maker *JWTMaker)CreateToken(id int, email string, isAdmin bool, duration time.Duration) (string, error) {
	claims, err := NewUserClaims(id, email, isAdmin, duration)
	if err != nil {
		return "", fmt.Errorf("error in creating claims")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", fmt.Errorf("error at signing the create token %v", err)
	}

	return tokenStr, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method: %v", token.Method)
		}
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		log.Printf("Token verification error: %v, secret key used for verification: %s", err, maker.secretKey)
		return nil, fmt.Errorf("invalid token signature")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}