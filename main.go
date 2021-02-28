package main

import (
	"log"
	"net/http"
	"os"
	"assignment-1/server"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	var rootPath string = "/exchange/v1/"

	http.HandleFunc(rootPath+"exchangehistory/", server.ExchangeHistoryHandler)
	http.HandleFunc(rootPath+"exchangeborder/", server.ExchangeBorderHandler)
	http.HandleFunc(rootPath+"diag/", server.DiagHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}