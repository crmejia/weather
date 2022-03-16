package weather_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"os"
	"testing"
	"weather"
)

func TestFormatURL(t *testing.T) {
	t.Parallel()
	location := "London"
	token := "dummy_token"
	want := "https://api.openweathermap.org/data/2.5/weather?q=London&appid=dummy_token"
	got := weather.FormatURL(location, token)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestParseJSON(t *testing.T) {
	t.Parallel()
	f, err := os.Open("testdata/london.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	want := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.17,
	}
	got := weather.ParseJSON(f)
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.001)) {
		t.Error(cmp.Diff(want, got))
	}
}
