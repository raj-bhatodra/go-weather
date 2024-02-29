package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Define your API key directly here
const OpenWeatherMapKey = "fd3c8caaa887ba999bdd97145cdb6289"

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func query(city string) (weatherData, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s", OpenWeatherMapKey, city))
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	// Print the response to the terminal
	fmt.Printf("Response: %+v\n", d)

	return d, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city, ok := vars["city"]
	if !ok {
		http.Error(w, "City not found in the URL", http.StatusBadRequest)
		return
	}

	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/weather/{city}", weatherHandler)

	http.Handle("/", r) // Register the router as the handler for incoming requests

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
