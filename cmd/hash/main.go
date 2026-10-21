package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := os.Getenv("PASSWORD_OWNER")

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		12,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(hash))
}