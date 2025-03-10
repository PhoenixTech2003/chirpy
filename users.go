package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/phoenixTech2003/chirpy/internal/auth"
	"github.com/phoenixTech2003/chirpy/internal/database"
)




func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Email  string `json:"email"`
		Password string `json:"password"`

	}

	type respBody struct {
		Id uuid.UUID `json:"id"`
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
	createUserParams := database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email: sql.NullString{String: params.Email, Valid: true},
	}
	
	user , err := cfg.dbQueries.CreateUser(r.Context(),createUserParams)
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

func (cfg *apiConfig) loginUser( w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds float64 `json:"expires_in_seconds"`
	}

	type responseParameters struct {
		Id string `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Email string `json:"email"`
		Token string `json:"token"`

	}

	decoder := json.NewDecoder(r.Body)
	requestParameters :=parameters{}
	err  := decoder.Decode(&requestParameters)
	if err != nil {
		log.Printf("An error occured while decoding the json %s",err)
		w.WriteHeader(500)
		return
	}
	userData, err := cfg.dbQueries.GetUserByEmail(r.Context(), sql.NullString{String: requestParameters.Email, Valid: true})
	if err != nil {
		log.Printf("Failed to login with error %s", err)
		w.WriteHeader(500)
		return
	} 

	err = auth.CheckPassword(requestParameters.Password, userData.HashedPassword)
	if err != nil {
		log.Printf("failed to authenticate %s", err)
		w.WriteHeader(401)
		return
	}
	fmt.Print(requestParameters.ExpiresInSeconds)
	if requestParameters.ExpiresInSeconds == 0 {
		requestParameters.ExpiresInSeconds = 3600
	}

	token, err := auth.MakeJWT(userData.ID, cfg.tokenSecret, time.Duration(requestParameters.ExpiresInSeconds))
	if err != nil {
		log.Printf("An error occured while generating jwt for user %s failed with error: %s",userData.ID.String(), err)
		w.WriteHeader(500)
		return
	}

	respBody := responseParameters{
		Id: userData.ID.String(),
		CreatedAt: userData.CreatedAt.Time,
		UpdatedAt: userData.CreatedAt.Time,
		Email: userData.Email.String,
		Token: token,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("an error occured while marshalling the response data %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.Write(dat)
}