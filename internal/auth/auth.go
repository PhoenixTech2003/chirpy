package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
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

func MakeJWT(userId uuid.UUID, tokenSecret string)(string, error){
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()) ,
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(3.6e+12))),
		Subject: userId.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedJWT, err :=token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("An error occured while signing your JWT: %s", err)
		return "" , err
	}

	return signedJWT , nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString,&claims,func(t *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil})
	if err != nil {
		log.Printf("An error occured while validating your jwt %s",err)
		return uuid.UUID{}, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("An error occured while getting the user id %s", err)
		return uuid.UUID{}, err
	}

	return uuid.MustParse(id) , nil

}

func GetBearerToken (headers http.Header) (string, error){
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == ""{
		log.Print("Authorization header not found")
		return "", errors.New("failed to get authorization header")
	}
	tokenString := strings.Split(authorizationHeader, " ")[1]
	return tokenString, nil
}

func MakeRefreshToken() (string, error){
	key := make([]byte,32)
	rand.Read(key)
	return hex.EncodeToString(key), nil
	
}