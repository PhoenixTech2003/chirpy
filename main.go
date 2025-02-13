package main

import (

	"net/http"
)


func main()  {
	mux := http.NewServeMux()
	imageMux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	imageMux.Handle("/assets", http.FileServer(http.Dir("./assets/logo.png")))
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	server.ListenAndServe()
	
}