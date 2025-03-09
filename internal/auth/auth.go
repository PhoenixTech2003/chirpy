package auth

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration)(string, error){
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()) ,
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userId.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256,claims)
	signedJWT, err :=token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("An error occured while signing your JWT: %s", err)
		return "" , err
	}

	return signedJWT , nil

}