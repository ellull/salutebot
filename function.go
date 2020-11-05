package salutebot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/ellull/salutebot/pkg/saluter"
)

var s *saluter.Saluter

func init() {
	var err error

	// get the filename from the FILENAME environment variable
	// defaulting to goodmorning.csv
	filename := "goodmorning.csv"
	if env := os.Getenv("FILENAME"); env != "" {
		filename = env
	}

	// create the Saluter
	s, err = saluter.NewFileSaluter(filename)
	if err != nil {
		log.Fatalln(err)
	}
}

type saluteRequest struct {
	URLs []string `json:"urls"`
}

type saluteResponse struct {
	Results []saluteResult `json:"results"`
}

type saluteResult struct {
	URL    string `json:"url"`
	Result string `json:"result"`
}

// SaluteHandler sends a salute to the Google Chat webhook specified in the URL environment
func SaluteHandler(w http.ResponseWriter, r *http.Request) {
	// only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// read the JSON request
	jsonRequest := saluteRequest{}
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil || len(jsonRequest.URLs) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// send the salutes
	jsonResponse := saluteResponse{Results: make([]saluteResult, len(jsonRequest.URLs))}
	for i, url := range jsonRequest.URLs {
		jsonResponse.Results[i] = saluteResult{URL: url, Result: result(s.Salute(url))}
	}

	// write the JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Printf("Error responding request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func result(err error) string {
	if err != nil {
		return err.Error()
	}
	return "ok"
}
