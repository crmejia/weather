package weather_test

import (
	"fmt"
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

func TestLocationFromArgsParsesLocationsCorrectly(t *testing.T) {
	testCases := []struct {
		input []string
		want  string
	}{
		{input: []string{"london"}, want: "london"},
		{input: []string{"london,", "uk"}, want: "london,%20uk"},
		{input: []string{"london", ",", "uk"}, want: "london,uk"},
		{input: []string{"santo", "domingo"}, want: "santo%20domingo"},
		{input: []string{"los", "angeles,", "us"}, want: "los%20angeles,%20us"},
	}

	for _, tc := range testCases {
		got, _ := weather.LocationFromArgs(tc.input)
		if tc.want != got {
			t.Errorf("want %q, got %q", tc.want, got)
		}
	}

}

func TestLocationFromArgsReturnErrorOnNoLocation(t *testing.T) {
	_, err := weather.LocationFromArgs([]string{})
	if err == nil {
		t.Error("want error on no location")
	}
}
func TestStringerOnConditions(t *testing.T) {
	cond := weather.Conditions{
		Summary:            "Drizzle",
		TemperatureCelsius: 7.2,
	}
	want := "Drizzle 7.2ºC"
	got := fmt.Sprint(cond)

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
