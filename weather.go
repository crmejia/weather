package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type responseJSON struct {
	Coordinates coordinates `json:"coord"`
	Weather     []weather   `json:"weather"`
	Main        main        `json:"main"`
}
type coordinates struct {
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}
type weather struct {
	Description string `json:"main""`
}
type main struct {
	Temp float32
}
type Conditions struct {
	Summary            string
	TemperatureCelsius float32
}

func FormatURL(location, token string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", location, token)
}

const kelvinToCelcius = 273.1500

func ParseJSON(r io.Reader) Conditions {
	decoder := json.NewDecoder(r)
	var resp responseJSON
	err := decoder.Decode(&resp)
	if err != nil {
		log.Println(err)
		return Conditions{}
	}
	var cond Conditions
	cond.Summary = resp.Weather[0].Description
	cond.TemperatureCelsius = resp.Main.Temp - kelvinToCelcius

	return cond
}

func Current(location, token string) (Conditions, error) {
	url := FormatURL(location, token)
	resp, err := http.Get(url)
	if err != nil {
		return Conditions{}, err
	}
	cond := ParseJSON(resp.Body)
	return cond, nil
}

func LocationFromArgs(input []string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("input location cannot be empty")
	}
	var output string
	unparsedComma := false
	for i, _ := range input {
		if i > 0 && input[i] != "," && !unparsedComma {
			output += "%20" + input[i]
		} else {
			if input[i] == "," {
				unparsedComma = true
			}
			output += input[i]
		}
	}
	return output, nil
}
