package auth

import (
	"errors"
	"fmt"
	"net/http"
	"resumegenerator/internal/database"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	signingMethod *jwt.SigningMethodHMAC
	signingKey    []byte
}

func New(signingKey string) *Auth {
	return &Auth{
		signingMethod: jwt.SigningMethodHS256,
		signingKey:    []byte(signingKey),
	}
}

func (a *Auth) DecodeAuthHeader(header http.Header) (*jwt.Token, error) {
	authorization := strings.Split(header.Get("authorization"), " ")
	if len(authorization) < 2 {
		return nil, errors.New("invalid authorization header")
	}
	tokenStr := authorization[1]
	return a.DecodeToken(tokenStr)
}

func (a *Auth) DecodeToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != a.signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(a.signingKey), nil
	})
}

func (a *Auth) GenToken(user *database.User) (string, error) {
	token := jwt.NewWithClaims(a.signingMethod, jwt.MapClaims{
		"userId": user.Id,
	})

	signed, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (a *Auth) TokenUserId(token *jwt.Token) (string, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userId, ok := claims["userId"].(string); ok {
			return userId, nil
		} else {
			return "", errors.New("invalid string")
		}
	} else {
		return "", errors.New("invalid claims")
	}
}
