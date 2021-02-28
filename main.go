package assignment_1

import (
	"log"
	"net/http"
	"os"
	"assignment-1/server"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT must be set")
	}

	var rootPath string = "/exchange/v1"

	http.HandleFunc(rootPath+"exchangehistory", server.ExchangeHistory)
	http.HandleFunc(rootPath+"exchangeborder", server.ExchangeBorder)
	http.HandleFunc(rootPath+"diag", server.Diag)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}