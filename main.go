package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type apiConfigData struct {
	OpenWeatherMapKey string `json:"OpenWeatherMapKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (*apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

func query(city string, apiKey string) (weatherData, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s", apiKey, city))
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city, ok := vars["city"]
	if !ok {
		http.Error(w, "City not found in the URL", http.StatusBadRequest)
		return
	}

	apiConfig, err := loadApiConfig(".env")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := query(city, apiConfig.OpenWeatherMapKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/hello", hello)
	r.HandleFunc("/weather/{city}", weatherHandler)

	http.Handle("/", r) // Register the router as the handler for incoming requests

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()

	fmt.Println("Server is listening on port 8080")
	select {}
}
