package weather

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	openweather_api_token = "OPENWEATHER_API_TOKEN"
	unit_usage            = "set the Unit: Celsius(default), Fahrenheit, kelvin"
	detailed_usage        = "gives more weather details"
	coord_usage           = "use coordinates to determine weather"
	shorthand             = " (shorthand)"
	help_intro            = `Weather is a Go library package and accompanying command-line tool that will (briefly) report the current weather conditions for a given location.
    weather <options> location`
)

func RunCLI() {
	var unit string
	flag.StringVar(&unit, "Unit", Celsius, unit_usage)
	flag.StringVar(&unit, "u", Celsius, unit_usage+shorthand)

	var detailed bool
	flag.BoolVar(&detailed, "detailed", false, detailed_usage)
	flag.BoolVar(&detailed, "d", false, detailed_usage)

	var lon float64
	var lat float64
	flag.Float64Var(&lon, "longitude", 0, coord_usage)
	flag.Float64Var(&lon, "lon", 0, coord_usage+shorthand)
	flag.Float64Var(&lat, "latitude", 0, coord_usage)
	flag.Float64Var(&lat, "lat", 0, coord_usage+shorthand)

	if len(os.Args) == 1 {
		fmt.Println(help_intro)
		flag.Usage()
		return
	}
	flag.Parse()

	token := os.Getenv(openweather_api_token)
	if token == "" {
		log.Fatal("please set the Open Weather API token(OPENWEATHER_API_TOKEN)")
	}

	location := ""
	var err error
	if lat == 0 && lon == 0 {
		location, err = LocationFromArgs(flag.Args())
		if err != nil {
			log.Fatal(err)
		}
	}

	clientConfig := ClientConfig{
		Token:          token,
		Unit:           unit,
		DetailedFormat: detailed,
		Location:       location,
		Lat:            lat,
		Lon:            lon,
	}
	client := NewClient(clientConfig)

	cond, err := client.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cond)
}
