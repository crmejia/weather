package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	CacheTime   time.Time
}

func (c *Conditions) Convert() {
	if c.Unit == "c" || c.Unit == Celsius {
		c.Temperature -= kelvinToCelcius
		c.TempMin -= kelvinToCelcius
		c.TempMax -= kelvinToCelcius
	} else if c.Unit == "f" || c.Unit == Fahrenheit {
		c.Temperature = (c.Temperature-kelvinToCelcius)*9/5 + 32
		c.TempMin = (c.Temperature-kelvinToCelcius)*9/5 + 32
		c.TempMax = (c.Temperature-kelvinToCelcius)*9/5 + 32
	}

}

func (c Conditions) String() string {
	var unit string
	if c.Unit != "" {
		unit = strings.ToUpper(c.Unit[0:1])
	} else {
		unit = "C"
	}
	if c.LongFormat {
		return fmt.Sprintf("%s %.1fº%s\n%s min %.1fº%s, max %.1fº%s", c.Name, c.Temperature, unit, c.Description, c.TempMin, unit, c.TempMax, unit)
	}
	return fmt.Sprintf("%s %.1fº%s", c.Summary, c.Temperature, unit)
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
	if len(resp.Weather) == 0 {
		log.Fatal("malformed response")
	}
	var cond Conditions
	cond.Name = resp.Name
	cond.Summary = resp.Weather[0].Main
	cond.Temperature = resp.Main.Temp
	cond.Longitude = resp.Coord.Lon
	cond.Latitude = resp.Coord.Lat
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
	Location       string
	Lat, Lon       float64
}

type Client struct {
	token          string
	HttpClient     http.Client
	Unit           string
	DetailedFormat bool
	Url            string
	CacheKey       string
}

func NewClient(config ClientConfig) Client {
	config.Unit = strings.ToLower(config.Unit)

	if config.Unit == "" || (config.Unit != Fahrenheit && config.Unit != "f" && config.Unit != Kelvin && config.Unit != "k") {
		config.Unit = Celsius
	}

	var url, cacheKey string
	if config.Lat != 0 && config.Lon != 0 {
		url = FormatURLByCoordinates(float32(config.Lat), float32(config.Lon), config.Token)
		latString := strconv.Itoa(int(config.Lat))
		lonString := strconv.Itoa(int(config.Lon))
		cacheKey = latString + lonString
	} else {
		url = FormatURLByLocation(config.Location, config.Token)
		cacheKey = config.Location
	}

	return Client{
		token:          config.Token,
		Unit:           config.Unit,
		HttpClient:     http.Client{},
		DetailedFormat: config.DetailedFormat,
		Url:            url,
		CacheKey:       cacheKey,
	}
}

func (c Client) Token() string {
	return c.token
}

func (c Client) Current() (Conditions, error) {
	cacheEntry := CacheRetrieve(c.CacheKey)
	cond := ParseCache(cacheEntry)
	if cond.Summary == "" {
		resp, err := c.HttpClient.Get(c.Url)
		if err != nil {
			return Conditions{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return Conditions{}, fmt.Errorf("received HTTP status %d on request", resp.StatusCode)
		}
		cond = ParseJSON(resp.Body)
		cond.Unit = c.Unit
		cond.LongFormat = c.DetailedFormat

		cond.Convert()
		//caching
		cond.CacheTime = time.Now()
		marshalCond, _ := json.Marshal(cond)
		CacheAdd(c.CacheKey, marshalCond)
	}
	return cond, nil
}

func CacheRetrieve(key string) []byte {
	tempDir := os.TempDir()
	f, err := os.Open(tempDir + key)
	if err != nil {
		return nil
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}

	return fileBytes
}

const cacheDuration = time.Minute * 15

func ParseCache(b []byte) Conditions {
	cond := Conditions{}
	err := json.Unmarshal(b, &cond)
	if err != nil {
		return Conditions{}
	}

	t := time.Now()
	if t.Sub(cond.CacheTime) > cacheDuration {
		return Conditions{}
	}

	return cond
}

func CacheAdd(key string, marshalCond []byte) error {
	f, err := os.Create(os.TempDir() + key)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(marshalCond)
	if err != nil {
		return err
	}
	return nil
}

func CacheDelete(key string) error {
	if key == "" {
		return errors.New("empty key")
	}

	err := os.Remove(os.TempDir() + key)
	if err != nil {
		return err
	}
	return nil
}

type responseJSON struct {
	Coord   coordinates
	Weather []weather
	Main    main
	Name    string
}
type coordinates struct {
	Lon float32
	Lat float32
}
type weather struct {
	Main        string
	Description string
}
type main struct {
	Temp    float32
	TempMin float32 `json:"temp_min"`
	TempMax float32 `json:"temp_max"`
}

const kelvinToCelcius = 273.1500
const (
	Celsius    = "celsius"
	Fahrenheit = "fahrenheit"
	Kelvin     = "kelvin"
)
