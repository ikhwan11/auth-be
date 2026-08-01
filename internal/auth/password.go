package auth

import "golang.org/x/crypto/bcrypt"

const bcryptCost = bcrypt.DefaultCost

// HashPassword mengubah password plain text menjadi hash bcrypt.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// CheckPassword membandingkan password plain text dengan hash bcrypt.
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
