package main

import (
	"encoding/json"
	"log"
	"net/http"
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
			Error: "Something went wrong",
		}
		dat, _ := json.Marshal(errorResp)
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
}