package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const apiKey = "8090242bfd48442288b94002230511&q"

type WeatherResponse struct {
	Location LocationData `json:"location"`
	Current  CurrentData  `json:"current"`
}

type LocationData struct {
	Name string `json:"name"`
}

type CurrentData struct {
	TempC float64 `json:"temp_c"`
}

func getWeatherForecast(apiKey, location string) (*WeatherResponse, error) {
	url := "https://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + location

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var response WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	location := mux.Vars(r)["location"]
	weatherData, err := getWeatherForecast(apiKey, location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{location}", GetWeatherHandler).Methods("GET")

	http.Handle("/", r)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
