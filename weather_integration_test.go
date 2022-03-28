//go:build integration
// +build integration

package weather_test

import (
	"os"
	"testing"
	"weather"
)

func TestConditionsIntegration(t *testing.T) {
	t.Parallel()
	token := os.Getenv("OPENWEATHER_API_TOKEN")
	if token == "" {
		t.Skip("Please set a valid API key in the environment variable OPENWEATHER_API_TOKEN")
	}
	clientConfig := weather.ClientConfig{
		Token: token,
	}
	wClient := weather.NewClient(clientConfig)
	url := weather.FormatURLByLocation("London", token)
	cond, err := wClient.Current(url)
	if err != nil {
		t.Fatal(err)
	}
	if cond.Summary == "" {
		t.Errorf("empty summary: %#v", cond)
	}
}
