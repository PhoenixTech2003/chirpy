package auth

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password),10)
	if err != nil {
		log.Printf("An error occured while hashing your password %s", err)
		return "", err
	}
	return string(hashedpassword) , nil	
}

func CheckPassword (password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("An error occured while verifying your password %s", err)
		return errors.New("invalid password")
	}

	return nil

}