# Assignment 1

A RESTapi project for the PROG2005 Cloud Technologies course. This README will provide instructions on installing GO and running the program.

**Clone the repository into a clean directory**

## Installation and run instructions - Windows, Linux and MacOS

**source:** https://golang.org/
**source:** https://restcountries.eu/
**source:** https://exchangeratesapi.io/

- Download and install the latest 1.15 version of GO from https://golang.org/dl/
    - To verify that GO is installed and u have the correct version, run the following command in your shell/terminal:
    ``` go version ```
    - Result should be:
    ``` go version go1.15.8 ```


- Navigate to the directory where you cloned the repository in your shell/terminal
    - Run:
    ``` go build . ```
    - Followed by:
    ``` go run main.go ```
    - The prompt will stand idle while it is running

- The program offers three services:
    - The output will allways be in JSON format
    - Mandatory input values are marked in the URL's provided with {<:value>}
    - Optional input values are marked in the URL's provided with {<value>}

    - You can find the exchange rate history of a given country in a given period of time:
    http://localhost:8080/exchange/v1/exchangehistory/<:country>/<:startyyyy>-<:startmm>-<:startdd>-<:endyyyy>-<:endmm>-<:enddd>
    
    - You can find bordering countries and information regarding their exchange rates:
    http://localhost:8080/exchange/v1/exchangeborder/<:countryname>?limit=<numberofcountries>

    - You can run the diagnostics tool to get the current status of the API's this API relies on:
    http://localhost:8080/exchange/v1/diag/
    
