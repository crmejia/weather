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
	Description string `json:"main"`
}
type main struct {
	Temp float32
}
type Conditions struct {
	Summary            string
	TemperatureCelsius float32
}

func (c Conditions) String() string {
	return fmt.Sprintf("%s %.1fºC", c.Summary, c.TemperatureCelsius)
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

func LocationFromArgs(input []string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("input location cannot be empty")
	}
	var output string
	unparsedComma := false
	for i := range input {
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

type Client struct {
	token      string
	HttpClient http.Client
}

func NewClient(token string) Client {
	return Client{
		token:      token,
		HttpClient: http.Client{},
	}
}

func (c Client) Token() string {
	return c.token
}

func (c Client) Current(url string) (Conditions, error) {
	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return Conditions{}, err
	}
	cond := ParseJSON(resp.Body)
	return cond, nil
}
