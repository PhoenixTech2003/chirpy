package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/phoenixTech2003/chirpy/internal/auth"
	"github.com/phoenixTech2003/chirpy/internal/database"
)

func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type respBody struct {
		Id         uuid.UUID `json:"id"`
		Created_At time.Time `json:"created_at"`
		Updated_At time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("An error occured while decoding the JSON: %s", err)
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
		Email:          sql.NullString{String: params.Email, Valid: true},
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		log.Printf("An error occured while creating the user %s", err)
	}

	response := respBody{
		Id:         user.ID,
		Created_At: user.CreatedAt.Time,
		Updated_At: user.UpdatedAt.Time,
		Email:      user.Email.String,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("An error occured while marshalling the json %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)

}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responseParameters struct {
		Id           string    `json:"id"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	requestParameters := parameters{}
	err := decoder.Decode(&requestParameters)
	if err != nil {
		log.Printf("An error occured while decoding the json %s", err)
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

	token, err := auth.MakeJWT(userData.ID, cfg.tokenSecret)
	if err != nil {
		log.Printf("An error occured while generating jwt for user %s failed with error: %s", userData.ID.String(), err)
		w.WriteHeader(500)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("An error occured while generating refresh token for user %s failed with error: %s", userData.ID.String(), err)
		w.WriteHeader(500)
		return
	}

	refreshtokenParams := database.CreateRefreshTokenParams{Token: refreshToken, UserID: uuid.NullUUID{UUID: userData.ID, Valid: true}}
	err = cfg.dbQueries.CreateRefreshToken(r.Context(), refreshtokenParams)
	if err != nil {
		log.Printf("An error occured while creating refresh token for user %s failed with error: %s", userData.ID.String(), err)
		w.WriteHeader(500)
		return
	}

	respBody := responseParameters{
		Id:           userData.ID.String(),
		CreatedAt:    userData.CreatedAt.Time,
		UpdatedAt:    userData.CreatedAt.Time,
		Email:        userData.Email.String,
		Token:        token,
		RefreshToken: refreshToken,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("an error occured while marshalling the response data %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}


func (cfg *apiConfig) updateEmailAndPassword(w http.ResponseWriter, r *http.Request ){
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type responseParameters struct {
		Id uuid.UUID `json:"id"`
		Email string `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("an error occured while wxtracting the access token %s", err)
		w.WriteHeader(401)
		return
	}

	userId , err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		log.Printf("an error occured while wxtracting the access token %s", err)
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)

	requestParams := parameters{}

	err = decoder.Decode(&requestParams)
	if err != nil {
		log.Printf("Failed to decode json with error %s", err)
		w.WriteHeader(500)
		return
	}

	hasedPassword,err := auth.HashPassword(requestParams.Password)
	if err != nil {
		log.Printf("an error occured while hashing your password, %s", err)
		w.WriteHeader(500)
		return
	}
	updatePasswordAndEmailParams := database.UpdatePasswordAndEmailParams{
		Email: sql.NullString{String: requestParams.Email, Valid: true},
		HashedPassword: hasedPassword,
		ID: userId,
	}
	userData, err := cfg.dbQueries.UpdatePasswordAndEmail(r.Context(), updatePasswordAndEmailParams)
	if err != nil {
		log.Printf("failed to update email and password")
		w.WriteHeader(500)
		return
	}

	responseBody := responseParameters{
		Id: userData.ID,
		Email: userData.Email.String,
		CreatedAt: userData.CreatedAt.Time,
		UpdatedAt: userData.UpdatedAt.Time,
	}

	dat , err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("an error occured while marshalling the json %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	w.Write(dat)

}