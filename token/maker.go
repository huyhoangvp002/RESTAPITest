package token

import "time"

type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, role string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
