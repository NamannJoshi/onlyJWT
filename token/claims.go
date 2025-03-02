package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool    `json:"isAdmin"`
	jwt.RegisteredClaims
}

func NewUserClaims(id int, email string, isAdmin bool, duration time.Duration) (*UserClaims, error) {
	iD, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error creating uuid")
	}
	return &UserClaims{
		ID: id,
		Email: email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: iD.String(),
			Subject: email,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}