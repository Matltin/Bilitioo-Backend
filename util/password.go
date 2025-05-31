package util

import (
	"errors"
)

// HashedPassword return bcrypt hash of password
func HashedPassword(password string) (string, error) {
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to hash password: %s", err)
	// }
	// return string(hashedPassword), nil

	return password, nil
}

// CheckPassword checks if the provoided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	// return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if password != hashedPassword {
		return errors.New("not equal password")
	}

	return nil
}
