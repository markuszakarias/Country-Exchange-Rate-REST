package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Global variables
var StartDate string
var EndDate string
var Limit int

// Struct for storing json data about a country
type Countries struct {
	Name     string `json:"name"`
	Currency []struct {
		Code string `json:"code"`
	} `json:"currencies"`
	Border []string `json:"borders"`
}

// Struct for storing json data about bordering countries
type BorderCountries struct {
	Rates map[string]interface{} `json:"rates"`
	Base  string                 `json:"base"`
}

// Struct for storing json data about exchange rates
type ExchangeRates struct {
	Rates   map[string]interface{} `json:"rates"`
	StartAt string                 `json:"start_at"`
	Base    string                 `json:"base"`
	EndAt   string                 `json:"end_at"`
}

// Struct for storing json data about border countries exchange rates
type BorderExchangeRates struct {
	Rates map[string]interface{} `json:"rates"`
	Names string                 `json:"name"`
	Base  string                 `json:"base"`
}

// Function handler for endpoint one
func ExchangeHistoryHandler(w http.ResponseWriter, r *http.Request) {
	reqURL := strings.Split(r.URL.String(), "/")
	reqDates := strings.Split(reqURL[5], "-")
	reqCountry := reqURL[4]

	if historyDataValidation(reqURL, reqDates, w) {
		countryURL := "https://restcountries.eu/rest/v2/name/" + reqCountry

		getCountries, err := http.Get(countryURL)
		if err != nil {
			fmt.Println("HTTP request failed with error #{err}")
		}

		var userCountry []Countries
		err = json.NewDecoder(getCountries.Body).Decode(&userCountry)
		if err != nil {
			fmt.Println("Error decoding json data - #{err}")
		}

		userCurrency := userCountry[0].Currency[0].Code

		exchangeURL := "https://api.exchangeratesapi.io/history?" +
			"start_at=" + StartDate +
			"&end_at=" + EndDate +
			"&symbols=" + userCurrency

		getExchange, err := http.Get(exchangeURL)
		if err != nil {
			fmt.Println("HTTP request failed with error #{err}")
		}

		var userBase ExchangeRates
		err = json.NewDecoder(getExchange.Body).Decode(&userBase)
		if err != nil {
			fmt.Println("Error decoding json data - #{err}")
		}

		if userBase.Base == "" {
			w.WriteHeader(http.StatusNoContent)
			fmt.Fprint(w, "Error status code: ", http.StatusNoContent)
			fmt.Fprintf(w, "The endpoint does not support EUR against EUR!\n")
			fmt.Fprintf(w, "Please try with another country and the following format:\n")
			fmt.Fprintf(w, "exchange/v1/exchangehistory/<countryname>/<start_yyyy>-<start_mm>-<start_dd>-<end_yyyy>-<end_mm>-<end_dd> \n")
		} else {
			jsonBytes, err := json.Marshal(userBase)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
		}
	}
}

func ExchangeBorderHandler(w http.ResponseWriter, r *http.Request) {
	reqURL := strings.Split(r.URL.String(), "/")
	reqLimit := strings.Split(r.URL.String(), "=")
	reqCountry := reqURL[4]

	if borderDataValidation(reqURL, w) {
		countryURL := "https://restcountries.eu/rest/v2/name/" + reqCountry

		getCountries, err := http.Get(countryURL)
		if err != nil {
			fmt.Println("HTTP request failed with error #{err}")
		}

		var userCountry []Countries
		err = json.NewDecoder(getCountries.Body).Decode(&userCountry)
		if err != nil {
			fmt.Println("Error decoding json data - #{err}")
		}

		borderLimitValidation(reqLimit, userCountry)

		var borderCountries BorderCountries
		borderCountries.Rates = make(map[string]interface{})

		borderCountries = generateBorderData(borderCountries, userCountry, w)

		output, err := json.Marshal(borderCountries)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(output)
	}
}

func historyDataValidation(reqURL []string, reqDates []string, w http.ResponseWriter) bool {
	if (len(reqURL) == 6) && (len(reqDates) == 6) {
		startYear, err := strconv.Atoi(reqDates[0])
		startMonth, err := strconv.Atoi(reqDates[1])
		startDay, err := strconv.Atoi(reqDates[2])
		endYear, err := strconv.Atoi(reqDates[3])
		endMonth, err := strconv.Atoi(reqDates[4])
		endDay, err := strconv.Atoi(reqDates[5])
		if err != nil {
			fmt.Println("Was not able to convert with error #{err}")
			return false
		}

		if (startYear > 1900 && startYear < 2022) && (endYear > 1900 && endYear < 2022) &&
			(startMonth < 13 && startMonth > 0) && (endMonth < 13 && endMonth > 0) &&
			(startDay < 32 && startDay > 0) && (endDay < 32 && endDay > 0) {

			StartDate = reqDates[0] + "-" + reqDates[1] + "-" + reqDates[2]
			EndDate = reqDates[3] + "-" + reqDates[4] + "-" + reqDates[5]

			return true
		}
	}

	fmt.Fprintf(w, "The endpoint format is incorrect, please use the following format:\n\n")
	fmt.Fprintf(w, "\texchange/v1/exchangehistory/<countryname>/<start_yyyy>-<start_mm>-<start_dd>-<end_yyyy>-<end_mm>-<end_dd>")
	return false
}

func borderDataValidation(reqURL []string, w http.ResponseWriter) bool {
	if len(reqURL) == 5 && reqURL[4] != "" {
		return true
	}
	fmt.Fprintf(w, "The endpoint format is incorrect, please use the following format:\n\n")
	fmt.Fprintf(w, "\texchange/v1/exchangeborder/<countryname>?limit=<numberofcountries>")
	return false
}

func borderLimitValidation(limit []string, reqCountry []Countries) {
	var err error

	if len(limit) > 1 {
		Limit, err = strconv.Atoi(limit[1])
		if err != nil {
			fmt.Printf("Failed to convert from string to int, with error #{err}")
		}
	} else {
		Limit = len(reqCountry[0].Border)
	}

	// If user input limit higher than actual border count, set to border count
	if Limit > len(reqCountry[0].Border) {
		Limit = len(reqCountry[0].Border)
	}
}

func generateBorderData(borderCountries BorderCountries, countries []Countries, w http.ResponseWriter) BorderCountries {
	for i := 0; i < Limit; i++ {

		borderCountriesURL := "https://restcountries.eu/rest/v2/alpha/" + countries[0].Border[i]
		getBorder, err := http.Get(borderCountriesURL)
		if err != nil {
			fmt.Printf("The HTTP request failed with error #{err}")
		}

		var borderCountry Countries
		err = json.NewDecoder(getBorder.Body).Decode(&borderCountry)
		if err != nil {
			fmt.Println("Error decoding json data - #{err}")
		}

		borderName := borderCountry.Name

		exchangeURL := "https://api.exchangeratesapi.io/latest?base=" + countries[0].Currency[0].Code
		getExchange, err := http.Get(exchangeURL)
		if err != nil {
			fmt.Println("HTTP request failed with error #{err}")
		}

		var userBase BorderExchangeRates
		err = json.NewDecoder(getExchange.Body).Decode(&userBase)
		if err != nil {
			fmt.Println("Error decoding json data - #{err}")
		}

		tempCountry := make(map[string]interface{})
		tempCountry["rate"] = userBase.Rates[borderCountry.Currency[0].Code]
		tempCountry["currency"] = borderCountry.Currency[0].Code

		borderCountries.Base = countries[0].Currency[0].Code

		if tempCountry["currency"] == borderCountries.Base {
			tempCountry["rate"] = 1
		}

		if tempCountry["rate"] != nil {
			borderCountries.Rates[borderName] = tempCountry
		}
	}
	return borderCountries
}



