package util

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	// This is a placeholder for actual password hashing logic.
	// In a real application, you would use a library like bcrypt to hash the password.
	return string(hashedPassword), nil // Replace with actual hashing logic
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
