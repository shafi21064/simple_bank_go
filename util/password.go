package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hassed the password
func HassedPassword(password string) (string, error) {
	hassedPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hass password")
	}
	return string(hassedPassWord), nil
}

// Check hassed password
func CheckPassword(password string, hassdPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hassdPassword), []byte(password))
}
