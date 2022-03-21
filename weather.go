package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type responseJSON struct {
	Coordinates coordinates `json:"coord"`
	Weather     []weather   `json:"weather"`
	Main        main        `json:"main"`
	Name        string
}
type coordinates struct {
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}
type weather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}
type main struct {
	Temp    float32
	TempMin float32 `json:"temp_min"`
	TempMax float32 `json:"temp_max"`
}
type Conditions struct {
	Name        string
	Summary     string
	Temperature float32
	Unit        string
	LongFormat  bool
	Longitude   float32
	Latitude    float32
	Description string
	TempMin     float32
	TempMax     float32
}

func (c Conditions) String() string {
	unit := strings.ToUpper(c.Unit[0:1])
	if c.LongFormat {
		return fmt.Sprintf("%s %.1fº%s\n%s min %.1fº%s, max %.1fº%s", c.Name, c.Temperature, unit, c.Description, c.TempMin, unit, c.TempMax, unit)
	}
	return fmt.Sprintf("%s %.1fº%s", c.Summary, c.Temperature, unit)
}

const kelvinToCelcius = 273.1500
const (
	CELCIUS    = "celcius"
	FAHRENHEIT = "fahrenheit"
	KELVIN     = "kelvin"
)

func (c *Conditions) Convert() {
	if c.Unit == "c" || c.Unit == CELCIUS {
		c.Temperature -= kelvinToCelcius
		c.TempMin -= kelvinToCelcius
		c.TempMax -= kelvinToCelcius
	} else if c.Unit == "f" || c.Unit == FAHRENHEIT {
		c.Temperature = (c.Temperature-kelvinToCelcius)*9/5 + 32
		c.TempMin = (c.Temperature-kelvinToCelcius)*9/5 + 32
		c.TempMax = (c.Temperature-kelvinToCelcius)*9/5 + 32
	}

}

func FormatURLByLocation(location, token string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", location, token)
}

func FormatURLByCoordinates(lat, lon float32, token string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%.2f&lon=%.2f&appid=%s", lat, lon, token)
}

func ParseJSON(r io.Reader) Conditions {
	decoder := json.NewDecoder(r)
	var resp responseJSON
	err := decoder.Decode(&resp)
	if err != nil {
		log.Println(err)
		return Conditions{}
	}
	var cond Conditions
	cond.Name = resp.Name
	cond.Summary = resp.Weather[0].Main
	cond.Temperature = resp.Main.Temp
	cond.Longitude = resp.Coordinates.Longitude
	cond.Latitude = resp.Coordinates.Latitude
	cond.Description = resp.Weather[0].Description
	cond.TempMin = resp.Main.TempMin
	cond.TempMax = resp.Main.TempMax

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

type ClientConfig struct {
	Token          string
	Unit           string
	DetailedFormat bool
}

type Client struct {
	token          string
	HttpClient     http.Client
	unit           string
	DetailedFormat bool
}

func NewClient(config ClientConfig) Client {
	return Client{
		token:          config.Token,
		unit:           config.Unit,
		HttpClient:     http.Client{},
		DetailedFormat: config.DetailedFormat,
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
	cond.Unit = c.unit
	cond.LongFormat = c.DetailedFormat
	cond.Convert()
	return cond, nil
}
