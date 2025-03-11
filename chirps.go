package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/phoenixTech2003/chirpy/internal/auth"
	"github.com/phoenixTech2003/chirpy/internal/database"
)

func (cfg *apiConfig) postChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorParameters struct {
		Body    string `json:"body"`
		User_id string `json:"user_id"`
		Error   string `json:"error"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed to extract token %s", err)
		w.WriteHeader(401)
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		log.Printf("Failed to extract token %s", err)
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("An error occured while decoding the json: %s", err)
		errorResp := errorParameters{
			Body:    "",
			User_id: "",
			Error:   "Something went wrong while recieving the request",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	createChirpParams := database.CreateChirpParams{
		UserID: uuid.NullUUID{UUID: userId, Valid: true},
		Body:   sql.NullString{String: params.Body, Valid: true},
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		log.Printf("An error occured while inserting chirp into database %s", err)
		errorResp := errorParameters{
			Body:    "",
			User_id: "",
			Error:   "Something went wrong while creating the chirp",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	postChirpResponse := parameters{

		Body: chirp.Body.String,
	}

	dat, err := json.Marshal(postChirpResponse)
	if err != nil {
		log.Printf("An error occured while inserting chirp into database %s", err)
		errorResp := errorParameters{
			Body:    "",
			User_id: "",
			Error:   "Something went wrong while creating the chirp",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Print("an error occured while getting bearer token")
		w.WriteHeader(401)
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		log.Print("an error occured while getting bearer token")
		w.WriteHeader(401)
		return
	}

	chirpId := r.PathValue("chirpId")

	deleteUserChirpParameters := database.DeleteUserChirpParams{
		ID:     uuid.MustParse(chirpId),
		UserID: uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = cfg.dbQueries.DeleteUserChirp(r.Context(), deleteUserChirpParameters)

	if err != nil {
		log.Printf("The record does not exist")
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)

}

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	if authorId != "" {
		chirps, err := cfg.dbQueries.GetAllChirpsByAuthor(r.Context(), uuid.NullUUID{UUID: uuid.MustParse(authorId), Valid: true})
		if err != nil {
			log.Printf("an error occured whiled getting chirps %s", err)
			w.WriteHeader(500)
			return
		}
	
		dat, err := json.Marshal(chirps)
		if err != nil {
			log.Printf("an error occured while marshalling the array of chirps")
			w.WriteHeader(500)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}

	

	
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("an error occured whiled getting chirps %s", err)
		w.WriteHeader(500)
		return
	}

	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("an error occured while marshalling the array of chirps")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)


}
