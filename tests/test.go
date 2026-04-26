package main

import (
	"fmt"
	"invoice-backend/Helper/Utils"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 🔐 Encrypt email
	encryptedEmail, err := Utils.EncryptForDB("admin@gmail.com")
	if err != nil {
		panic(err)
	}

	// 🔑 Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("Encrypted Email:")
	fmt.Println(encryptedEmail)

	fmt.Println("\nHashed Password:")
	fmt.Println(string(hashedPassword))
}