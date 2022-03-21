# Weather
Weather is a Go library package and accompanying command-line tool that will (briefly) report the current weather conditions for a given location. 

## Installation
1. Install by running `go install https://github.com/crmejia/weather/cmd/weather@latest`
2. To get started you'll need to sign up in to [OpenWeather]((https://openweathermap.org/) and generate an API Token. Once the API Token is created
   you'll need to export it into your environment like so: `export OPENWEATHER_API_TOKEN=<your token>`.

## Usage
`weather <options> location`

The options are 
```bash
  -d	gives more weather details
  -detailed
    	gives more weather details
  -lat float
    	use coordinates to determine weather (shorthand)
  -latitude float
    	use coordinates to determine weather
  -lon float
    	use coordinates to determine weather (shorthand)
  -longitude float
    	use coordinates to determine weather
  -u string
    	set the unit: CELCIUS(default), FAHRENHEIT, kelvin (shorthand) (default "celcius")
  -unit string
    	set the unit: CELCIUS(default), FAHRENHEIT, kelvin (default "celcius")
 ```
For example, 
```
weather London, UK
Cloudy 15.2ºC
```