package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)


func validator (w http.ResponseWriter, r *http.Request) {
	
	type parameters struct {
		Body string `json:"body"`
	}

	type errRespBody struct {
		Error string `json:"error"`
	}

	type returnVals struct {
		Valid string `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorMessage := errRespBody{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(errorMessage)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	lengthOfChirp := len(params.Body)
	if lengthOfChirp > 140 {
		errorMessage := errRespBody{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errorMessage)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	convertedMessageBody:= string(params.Body)
	messageBodySlice := strings.Split(convertedMessageBody," ")
	for _ , profaneWord := range profane{
		for i , word := range messageBodySlice {
			if profaneWord == word {
				messageBodySlice[i] = "****"
				
			}
		}
	}
	respBody := returnVals{
		Valid: strings.Join(messageBodySlice," "),
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
	
}

