package server

import (
	"net/http"
	"strings"
)

func ExchangeHistory(w http.ResponseWriter, r *http.Request) {

	url := strings.Split(r.URL.String(), "/")
	dates := strings.Split(url[5], "-")
	country := url[4]


}

func ExchangeBorder(w http.ResponseWriter, r *http.Request) {

}

func Diag(w http.ResponseWriter, r *http.Request) {

}
