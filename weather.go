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
	Summary     string
	Temperature float32
	Unit        string
}

func (c Conditions) String() string {
	unit := strings.ToUpper(c.Unit[0:1])
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
	} else if c.Unit == "f" || c.Unit == FAHRENHEIT {
		c.Temperature = (c.Temperature-kelvinToCelcius)*9/5 + 32
	}

}

func FormatURL(location, token string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", location, token)
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
	cond.Summary = resp.Weather[0].Description
	cond.Temperature = resp.Main.Temp

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
	Token string
	Unit  string
}

type Client struct {
	token      string
	HttpClient http.Client
	unit       string
}

func NewClient(config ClientConfig) Client {
	return Client{
		token:      config.Token,
		unit:       config.Unit,
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
	cond.Unit = c.unit
	cond.Convert()
	return cond, nil
}
