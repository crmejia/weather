package weather

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	openweather_api_token = "OPENWEATHER_API_TOKEN"
	unit_usage            = "set the unit: CELCIUS(default), FAHRENHEIT, kelvin"
	long_usage            = "gives more weather details"
	shorthand             = " (shorthand)"
)

func RunCLI() {

	var unit string
	flag.StringVar(&unit, "unit", CELCIUS, unit_usage)
	flag.StringVar(&unit, "u", CELCIUS, unit_usage+shorthand)

	var long bool
	flag.BoolVar(&long, "long", false, long_usage)
	flag.BoolVar(&long, "l", false, long_usage)
	flag.Parse()

	unit = strings.ToLower(unit)
	if unit != FAHRENHEIT && unit != "f" && unit != KELVIN && unit != "k" {
		unit = CELCIUS
	}

	location, err := LocationFromArgs(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv(openweather_api_token)
	if token == "" {
		log.Fatal("please set the Open Weather API token(OPENWEATHER_API_TOKEN)")
	}
	clientConfig := ClientConfig{
		Token:      token,
		Unit:       unit,
		LongFormat: long,
	}
	client := NewClient(clientConfig)
	url := FormatURL(location, token)
	cond, err := client.Current(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cond)
}
