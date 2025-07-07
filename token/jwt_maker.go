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
	fmt.Println("[DEBUG] NewJWTMaker called")
	fmt.Println("[DEBUG] Provided secretKey:", secretKey)

	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// fmt.Println("[DEBUG] CreateToken called")
	// fmt.Println("[DEBUG] Username:", username)
	// fmt.Println("[DEBUG] Duration:", duration)
	fmt.Println("[DEBUG] Using secretKey:", maker.secretKey)

	payload, err := NewPayload(username, duration)
	if err != nil {
		fmt.Println("[ERROR] Failed to create payload:", err)
		return "", err
	}

	fmt.Printf("[DEBUG] Created payload: %+v\n", payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		fmt.Println("[ERROR] Failed to sign token:", err)
		return "", err
	}

	fmt.Println("[DEBUG] Generated JWT Token:", tokenString)
	return tokenString, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	fmt.Println("[DEBUG] VerifyToken called")
	fmt.Println("[DEBUG] Incoming token:", token)
	fmt.Println("[DEBUG] Using secretKey:", maker.secretKey)

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		fmt.Println("[DEBUG] Inside keyFunc")
		fmt.Printf("[DEBUG] Token header: %+v\n", t.Header)

		alg := t.Header["alg"]
		fmt.Println("[DEBUG] Algorithm in header:", alg)

		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println("[ERROR] Unexpected signing method:", t.Method.Alg())
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		fmt.Println("[ERROR] jwt.ParseWithClaims failed:", err)
		if verr, ok := err.(*jwt.ValidationError); ok {
			fmt.Println("[DEBUG] ValidationError details:", verr)
			if (verr.Errors & jwt.ValidationErrorExpired) != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		fmt.Println("[ERROR] Failed to cast claims to Payload")
		return nil, ErrInvalidToken
	}

	fmt.Printf("[DEBUG] Successfully verified payload: %+v\n", payload)
	return payload, nil
}
