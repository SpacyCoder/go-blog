package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spacycoder/go-blog/accountservice/dbclient"
	"github.com/spacycoder/go-blog/accountservice/model"
	cb "github.com/spacycoder/go-blog/common/circuitbreaker"
)

var isHealthy = true
var DBClient dbclient.IBoltClient

var client = &http.Client{}

var fallbackQuote = model.Quote{
	Language: "en",
	ServedBy: "circuit-breaker",
	Text:     "May the source be with you. Always."}

func init() {
	var transport http.RoundTripper = &http.Transport{DisableKeepAlives: true}
	client.Transport = transport
}

func GetAccount(w http.ResponseWriter, r *http.Request) {

	accountID := mux.Vars(r)["accountId"]
	account, err := DBClient.QueryAccount(accountID)
	account.ServedBy = getIP()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	quote := getQuote()
	if err == nil {
		account.Quote = quote
	}
	data, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getQuote() model.Quote {

	body, err := cb.CallUsingCircuitBreaker("quotes-service", "http://quotes-service:8080/api/quote?strength=13", "GET")
	if err == nil {
		quote := model.Quote{}
		json.Unmarshal(body, &quote)
		return quote
	} else {
		return fallbackQuote
	}
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	panic("Unable to determine local IP address (non loopback). Exiting.")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbUp := DBClient.Check()
	if dbUp && isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccessible"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)

	}
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

func SetHealthyState(w http.ResponseWriter, r *http.Request) {

	var state, err = strconv.ParseBool(mux.Vars(r)["state"])

	if err != nil {
		fmt.Println("Invalid request to SetHealthyState, allowed values are true or false")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHealthy = state
	w.WriteHeader(http.StatusOK)
}
