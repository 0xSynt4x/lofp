package auth

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword checks a password against a bcrypt hash.
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength checks that a password meets minimum requirements.
func ValidatePasswordStrength(password string) error {
	if len(password) < 10 {
		return fmt.Errorf("password must be at least 10 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	lower := strings.ToLower(password)
	for _, banned := range commonPasswords {
		if lower == banned {
			return fmt.Errorf("that password is too common, please choose a different one")
		}
	}

	return nil
}

var commonPasswords = []string{
	"password1!", "password12", "password123", "qwerty1234", "letmein123",
	"1234567890", "123456789!", "abcdefghij", "qwertyuiop",
	"iloveyou12", "trustno1!!", "welcome123", "monkey12345",
	"master1234", "dragon1234", "football12", "shadow1234",
	"sunshine12", "princess12", "legends123", "legends1234",
	"legendsofp", "futurepast", "shattered1", "andor12345",
}
