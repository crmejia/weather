package weather_test

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"weather"
)

func TestFormatURLByLocation(t *testing.T) {
	t.Parallel()
	location := "London"
	token := "dummy_token"
	want := "https://api.openweathermap.org/data/2.5/weather?q=London&appid=dummy_token"
	got := weather.FormatURLByLocation(location, token)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestFormatURLByCoordinates(t *testing.T) {
	t.Parallel()
	var lon float32 = 33
	var lat float32 = 44
	token := "dummy_token"
	want := "https://api.openweathermap.org/data/2.5/weather?lat=44.00&lon=33.00&appid=dummy_token"
	got := weather.FormatURLByCoordinates(lat, lon, token)
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
		Name:        "London",
		Summary:     "Drizzle",
		Temperature: 7.17,
		Unit:        weather.Celsius,
		Longitude:   -0.13,
		Latitude:    51.51,
		Description: "light intensity drizzle",
		TempMin:     6,
		TempMax:     8,
	}
	got := weather.ParseJSON(f)
	got.Unit = weather.Celsius
	got.Convert()
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.001)) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestLocationFromArgsParsesLocationsCorrectly(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	_, err := weather.LocationFromArgs([]string{})
	if err == nil {
		t.Error("want error on no location")
	}
}
func TestStringerOnConditions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		cond weather.Conditions
		want string
	}{
		{cond: weather.Conditions{Summary: "Drizzle", Temperature: 7.2, Unit: "c"}, want: "Drizzle 7.2ºC"},
		{cond: weather.Conditions{Summary: "Drizzle", Temperature: 7.2, Unit: weather.Celsius}, want: "Drizzle 7.2ºC"},
		{cond: weather.Conditions{Summary: "Drizzle", Temperature: 7.2, Unit: weather.Fahrenheit}, want: "Drizzle 7.2ºF"},
		{cond: weather.Conditions{Summary: "Drizzle", Temperature: 7.2, Unit: weather.Kelvin}, want: "Drizzle 7.2ºK"},
	}

	for _, tc := range testCases {
		got := fmt.Sprint(tc.cond)
		if tc.want != got {
			t.Errorf("want %q, got %q", tc.want, got)
		}
	}
}

func TestStringerOnConditionsWithLongFormat(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		cond weather.Conditions
		want string
	}{
		{cond: weather.Conditions{
			LongFormat:  true,
			Name:        "Santo Domingo",
			Description: "light intensity drizzle",
			Temperature: 27.2,
			Longitude:   -23.4,
			Latitude:    43.51,
			TempMin:     22,
			TempMax:     30.3,
			Unit:        "c"},
			want: "Santo Domingo 27.2ºC\nlight intensity drizzle min 22.0ºC, max 30.3ºC"},
	}

	for _, tc := range testCases {
		got := fmt.Sprint(tc.cond)
		if tc.want != got {
			t.Errorf("want %q, got %q", tc.want, got)
		}
	}
}

var dummyClientConfig = weather.ClientConfig{
	Token: "dummy_token",
	Unit:  "c",
}

func TestConditionsConvertToAppropiateUnit(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		cond weather.Conditions
		want float32
	}{
		{cond: weather.Conditions{Temperature: 300, Unit: "c"}, want: 26.85},
		{cond: weather.Conditions{Temperature: 300, Unit: weather.Celsius}, want: 26.85},
		{cond: weather.Conditions{Temperature: 300, Unit: "k"}, want: 300},
		{cond: weather.Conditions{Temperature: 300, Unit: weather.Kelvin}, want: 300},
		{cond: weather.Conditions{Temperature: 300, Unit: "f"}, want: 80.33},
		{cond: weather.Conditions{Temperature: 300, Unit: weather.Fahrenheit}, want: 80.33},
	}

	for _, tc := range testCases {
		tc.cond.Convert()
		got := tc.cond.Temperature
		if !cmp.Equal(tc.want, got, cmpopts.EquateApprox(0, 0.001)) {
			t.Errorf("want %.1f, got %.1f", tc.want, got)
		}
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	c := weather.NewClient(dummyClientConfig)
	want := dummyClientConfig.Token
	got := c.Token()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestNewClientSetsDefaultUnitAsCelsius(t *testing.T) {
	t.Parallel()
	testConfigs := []weather.ClientConfig{
		weather.ClientConfig{Token: "dummy"},
		weather.ClientConfig{Token: "dummy", Unit: "gibberish"},
	}
	for _, testConfig := range testConfigs {
		client := weather.NewClient(testConfig)
		want := weather.Celsius
		got := client.Unit

		if want != got {
			t.Errorf("want unit %s, got %s", want, got)
		}
	}
}

func TestClient_Current(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(LoadLondonJSON))

	wclient := weather.NewClient(dummyClientConfig)
	wclient.HttpClient = *ts.Client()
	cond, _ := wclient.Current(ts.URL)

	want := "Drizzle 7.2ºC"
	got := cond.String()

	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func LoadLondonJSON(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("testdata/london.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fileBytes, _ := ioutil.ReadAll(f)

	w.Write(fileBytes)
}

func TestParseCacheSkipsStaleConditions(t *testing.T) {
	t.Parallel()
	marshalCond, _ := json.Marshal(weather.Conditions{
		CacheTime: time.Date(1970, time.January, 1, 0, 0, 0, 0, &time.Location{}),
	})
	want := weather.Conditions{}
	got := weather.ParseCache(marshalCond)
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.001)) {
		t.Error("stale cache shouldn't be parsed ")
	}
}

func TestParseCacheFreshConditions(t *testing.T) {
	t.Parallel()
	timeNow := time.Now()
	marshalCond, _ := json.Marshal(weather.Conditions{
		Name:      "test",
		CacheTime: timeNow,
	})
	want := weather.Conditions{
		Name:      "test",
		CacheTime: timeNow,
	}
	got := weather.ParseCache(marshalCond)
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.001)) {
		t.Error("stale cache shouldn't be parsed ")
	}
}

//
//func TestCacheRetrieve(t *testing.T) {
//	reader := "reader" //TODO json conditions
//
//}

//func TestCacheRetrieveReturnsErrorOnCacheMiss(t *testing.T) {
//	t.Parallel()
//	_, err := weather.CacheRetrieve("dummy")
//	if err == nil {
//		t.Errorf("expected retrieve to fail on non-existent item")
//	}
//}
//
//func TestCacheAddSavesToTmpFile(t *testing.T) {
//	t.Parallel()
//	key := "test"
//	cond := weather.Conditions{Name: "test"}
//	weather.CacheAdd(key, cond)
//
//	_, err := os.Stat(os.TempDir() + key)
//	if err == os.ErrNotExist {
//		t.Errorf("want conditions to be cached")
//	}
//}

//func TestCacheIsRetrieved
//func TestCacheRemovesStale
