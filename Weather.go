package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func fetchWeatherData() (map[string]interface{}, error) {
	apiKey := "603604d84c5de6f005987b06c9f94d9d"
	lat, lon := "49.319552", "8.5498581"
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, apiKey)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
