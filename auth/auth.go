package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = []byte("project kanban board")

type JWTClaim struct {
	Username string `json:"username"`
	Role     uint8  `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(username string, role uint8) (tokenString string, err error) {
	expTime := time.Now().Add(1 * 24 * time.Hour)
	claims := &JWTClaim{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(key)

	return
}

func ValidateToken(signedToken string) (role uint8, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)

	if !ok {
		err = errors.New("could not parse claim from token")

		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token expired")

		return
	}

	role = claims.Role

	return
}
