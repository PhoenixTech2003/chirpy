package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/phoenixTech2003/chirpy/internal/auth"
)

func (cfg *apiConfig) upgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	type upgradeToChirpyRedData struct {
		UserId string `json:"user_id"`
	}

	type upgradeToChirpyRedParameters struct {
		Event string                 `json:"event"`
		Data  upgradeToChirpyRedData `json:"data"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil{
		log.Printf("Failed to get api key, %s", err)
		w.WriteHeader(401)
		return
	}

	if tokenString != cfg.polkaApiKey {
		log.Print("api keys do not match")
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	upgradeToChirpyRedParams := upgradeToChirpyRedParameters{}
	err = decoder.Decode(&upgradeToChirpyRedParams)
	if err != nil {
		log.Printf("An error occured while decoding the data %s", err)
		w.WriteHeader(500)
		return
	}

	if upgradeToChirpyRedParams.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	userId := upgradeToChirpyRedParams.Data.UserId
	err = cfg.dbQueries.UpgradeToChirpyRed(r.Context(), uuid.MustParse(userId))
	if err != nil {
		log.Printf("An error occured while upgrading user to chirpy red %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}
