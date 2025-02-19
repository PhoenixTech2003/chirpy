package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)




func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Email  string `json:"email"`
	}

	type respBody struct {
		Id uuid.NullUUID `json:"id"`
		Created_At time.Time `json:"created_at"`
		Updated_At time.Time `json:"updated_at"`
		Email  string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("An error occured while decoding the JSON: %s", err )
		w.WriteHeader(500)
		return
	}
	
	user , err := cfg.dbQueries.CreateUser(r.Context(), sql.NullString{String: params.Email, Valid: true})
	if err != nil {
		log.Printf("An error occured while creating the user %s", err)
	}
	response := respBody{
		Id: user.ID,
		Created_At: user.CreatedAt.Time,
		Updated_At: user.UpdatedAt.Time,
		Email: user.Email.String,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("An error occured while marshalling the json %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	w.Write(dat)

}