package token

import "time"

type Maker interface {
	// CreateToken creates a new token for sepcific and duration
	CreateToken(UserID int64, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the tokern is valid or not
	VerifyToken(token string) (*Payload, error)
}
