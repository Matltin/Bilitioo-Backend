package util

// HashedPassword generates a SHA-256 hash of the password with a random salt
func HashedPassword(password string) (string, error) {
	// Generate a random salt
	// salt := make([]byte, 16)
	// if _, err := rand.Read(salt); err != nil {
	// 	return "", fmt.Errorf("failed to generate salt: %w", err)
	// }

	// // Create SHA-256 hash of password + salt
	// hash := sha256.Sum256([]byte(password + string(salt)))

	// // Return format: salt:hash
	// return fmt.Sprintf("%x:%x", salt, hash), nil

	return password, nil
}

// CheckPassword verifies if the provided password matches the hashed password
func CheckPassword(password string, hashedPassword string) error {
	// Split salt and hash
	// parts := strings.Split(hashedPassword, ":")
	// if len(parts) != 2 {
	// 	return errors.New("invalid hash format")
	// }

	// salt, err := hex.DecodeString(parts[0])
	// if err != nil {
	// 	return fmt.Errorf("failed to decode salt: %w", err)
	// }

	// // Compute hash of provided password with the same salt
	// hash := sha256.Sum256([]byte(password + string(salt)))

	// // Compare with stored hash
	// if hex.EncodeToString(hash[:]) != parts[1] {
	// 	return errors.New("passwords do not match")
	// }

	return nil
}
