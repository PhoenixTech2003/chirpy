package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/phoenixTech2003/chirpy/internal/auth"
)

func (cfg *apiConfig) postRefreshToken(w http.ResponseWriter, r *http.Request) {
	type respoBody struct {
		Body string `json:"body"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("an error occured while getting the bearer token from the header")
		w.WriteHeader(401)
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("an error occured while getting the bearer token from the header %s", err)
		w.WriteHeader(401)
		return

	}
	if time.Now().After(refreshToken.ExpiresAt.Time) {
		log.Printf("Your refresh token has expired please login to get another")
		w.WriteHeader(401)
		return
	}

	newTokenString, err := auth.MakeJWT(refreshToken.UserID.UUID, cfg.tokenSecret)
	if err != nil {
		log.Printf("an error occured while generating your jwt %s", err)
		w.WriteHeader(401)
		return
	}

	responseBody := respoBody{
		Body: newTokenString,
	}

	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("an error occured while marshalling json %s", err)
		w.WriteHeader(401)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) postRevokeToken(w http.ResponseWriter, r *http.Request){
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("an error occured while extracting your bearer token %s", err)
		w.WriteHeader(500)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(),tokenString)
	if err != nil {
		log.Printf("an error occured while revoking the token string %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
	
}
