package main

import (

	"net/http"
)


func main()  {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app",http.FileServer(http.Dir("."))),)
	mux.Handle("/app/assets", http.FileServer(http.Dir("./assets/logo.png")) ,)
	mux.HandleFunc("/healthz", func (response http.ResponseWriter, request *http.Request){
		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		response.WriteHeader(200)
		response.Write([]byte("OK"))
		
	})
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	server.ListenAndServe()
	
}