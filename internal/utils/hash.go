package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcryptCost,
	)
	if err != nil {
		return "", fmt.Errorf("gagal melakukan hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, passwordHash string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	); err != nil {
		return fmt.Errorf("password tidak sesuai: %w", err)
	}

	return nil
}