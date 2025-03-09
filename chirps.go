package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/phoenixTech2003/chirpy/internal/database"
)


func (cfg *apiConfig) postChirps(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Body string `json:"body"`
		User_id string `json:"user_id"`
	}
	type errorParameters struct {
		Body string `json:"body"`
		User_id string `json:"user_id"`
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("An error occured while decoding the json: %s", err)
		errorResp := errorParameters{
			Body: "",
			User_id: "",
			Error: "Something went wrong while recieving the request",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	createChirpParams := database.CreateChirpParams{
		UserID: uuid.NullUUID{Valid: true, UUID: uuid.MustParse(params.User_id)},
		Body: sql.NullString{String: params.Body, Valid: true},
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(),createChirpParams)
	if err != nil{
		log.Printf("An error occured while inserting chirp into database %s",err)
		errorResp := errorParameters{
			Body: "",
			User_id: "",
			Error: "Something went wrong while creating the chirp",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	postChirpResponse := parameters{
		User_id: chirp.ID.URN(),
		Body: chirp.Body.String,
	}

	dat , err := json.Marshal(postChirpResponse)
	if err != nil{
		log.Printf("An error occured while inserting chirp into database %s",err)
		errorResp := errorParameters{
			Body: "",
			User_id: "",
			Error: "Something went wrong while creating the chirp",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	w.Write(dat)

	
}