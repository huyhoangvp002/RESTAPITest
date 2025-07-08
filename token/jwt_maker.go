package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {

	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, role, duration)
	if err != nil {

		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {

		return "", err
	}

	return tokenString, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(t *jwt.Token) (interface{}, error) {

		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {

			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			if (verr.Errors & jwt.ValidationErrorExpired) != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
