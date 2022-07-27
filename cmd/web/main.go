package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeSiteHandler)
	mux.HandleFunc("/showguides", ShowGuidesHandler)
	mux.HandleFunc("/createguides", CreateGuideHandler)

	log.Println("Starting Server on Port :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
