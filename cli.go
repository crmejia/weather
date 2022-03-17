package weather

import (
	"fmt"
	"log"
	"os"
)

const openweather_api_token = "OPENWEATHER_API_TOKEN"

func RunCLI() {
	location, err := LocationFromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	token := os.Getenv(openweather_api_token)
	if token == "" {
		log.Fatal("please set the Open Weather API token(OPENWEATHER_API_TOKEN)")
	}
	client := NewClient(token)
	url := FormatURL(location, token)
	cond, err := client.Current(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cond)
}
