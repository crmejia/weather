package weather

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const openweather_api_token = "OPENWEATHER_API_TOKEN"

func RunCLI() {
	unit := flag.String("unit", CELCIUS, "set the unit: CELCIUS(default), FAHRENHEIT, kelvin")
	flag.Parse()

	if *unit != FAHRENHEIT && *unit != "f" && *unit != KELVIN && *unit != "k" {
		*unit = CELCIUS
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
		Token: token,
		Unit:  *unit,
	}
	client := NewClient(clientConfig)
	url := FormatURL(location, token)
	cond, err := client.Current(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cond)
}
