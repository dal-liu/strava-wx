package weather

import (
	"math"
	"testing"
)

func TestIsDay(t *testing.T) {
	tests := map[string]struct {
		input  weatherData
		result bool
	}{
		"day": {
			input:  weatherData{Dt: 2, Sunrise: 1, Sunset: 3},
			result: true,
		},
		"night": {
			input:  weatherData{Dt: 4, Sunrise: 1, Sunset: 3},
			result: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got, expected := test.input.isDay(), test.result; got != expected {
				t.Fatalf("isDay() got %t, expected %t", got, expected)
			}
		})
	}
}

func TestGetCondition(t *testing.T) {
	tests := map[string]struct {
		input     weatherData
		result    string
		resultErr string
	}{
		"no weather condition received": {
			input:     weatherData{},
			resultErr: "No weather condition received",
		},
		"unknown weather condition": {
			input:     weatherData{Weather: []weatherCondition{{0}}},
			resultErr: "Unknown weather condition",
		},
		"cloudy": {
			input:  weatherData{Weather: []weatherCondition{{804}}},
			result: "☁️ Cloudy",
		},
		"sunny": {
			input:  weatherData{Dt: 2, Sunrise: 1, Sunset: 3, Weather: []weatherCondition{{800}}},
			result: "☀️ Sunny",
		},
		"clear": {
			input:  weatherData{Dt: 4, Sunrise: 1, Sunset: 3, Weather: []weatherCondition{{800}}},
			result: "🌙 Clear",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cond, err := test.input.getCondition()

			if test.resultErr != "" {
				if err == nil {
					t.Fatalf("getCondition() got nil error, expected %s", test.resultErr)
				}
				if err.Error() != test.resultErr {
					t.Fatalf("getCondition() got error %s, expected %s", err.Error(), test.resultErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("getCondition() got error %s, expected %s", err.Error(), test.result)
			}
			if cond != test.result {
				t.Fatalf("getCondition() got %s, expected %s", cond, test.result)
			}
		})
	}
}

func TestGetWindDirection(t *testing.T) {
	tests := map[string]struct {
		input  weatherData
		result string
	}{
		"0": {
			input:  weatherData{Wind_deg: 0},
			result: "N",
		},
		"360": {
			input:  weatherData{Wind_deg: 360},
			result: "N",
		},
		"310": {
			input:  weatherData{Wind_deg: 310},
			result: "NW",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got, expected := test.input.getWindDirection(), test.result; got != expected {
				t.Fatalf("getWindDirection() got %s, expected %s", got, expected)
			}
		})
	}
}

func TestGetPrecipitation(t *testing.T) {
	tests := map[string]struct {
		input  weatherData
		result float64
	}{
		"0.0": {
			input:  weatherData{},
			result: 0.0,
		},
		"1.0": {
			input:  weatherData{Rain: weatherPrecipitation{12.7}, Snow: weatherPrecipitation{12.7}},
			result: 1.0,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got, expected := test.input.getPrecipitation(), test.result; math.Abs(got-expected) >= epsilon {
				t.Fatalf("getPrecipitation() got %f, expected %f", got, expected)
			}
		})
	}
}

func TestGetDescription(t *testing.T) {
	tests := map[string]struct {
		input     weatherResponse
		result    string
		resultErr string
	}{
		"no weather data": {
			input:     weatherResponse{},
			resultErr: "No weather data received",
		},
		"no weather condition received": {
			input:     weatherResponse{[]weatherData{{}}},
			resultErr: "No weather condition received",
		},
		"rain": {
			input:  weatherResponse{[]weatherData{{Temp: 50.0, Feels_like: 50.0, Humidity: 50, Wind_speed: 15.0, Wind_deg: 180, Wind_gust: 15.0, Weather: []weatherCondition{{500}}, Rain: weatherPrecipitation{0.5}}}},
			result: "🌧️ Rain, 50°F, Feels like 50°F, Humidity 50%, Wind 15mph with 15mph gusts from S, Precipitation 0.02 in/hr",
		},
		"rounding": {
			input:  weatherResponse{[]weatherData{{Temp: 67.8, Feels_like: 70.4, Humidity: 80, Wind_speed: 4.5, Wind_deg: 0, Weather: []weatherCondition{{802}}}}},
			result: "☁️ Partly cloudy, 68°F, Feels like 70°F, Humidity 80%, Wind 5mph from N",
		},
		"no wind": {
			input:  weatherResponse{[]weatherData{{Temp: 32.0, Feels_like: 20.0, Humidity: 41, Wind_speed: 0.0, Weather: []weatherCondition{{600}}, Snow: weatherPrecipitation{1.6}}}},
			result: "🌨️ Snow, 32°F, Feels like 20°F, Humidity 41%, Wind 0mph, Precipitation 0.06 in/hr",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			desc, err := test.input.getDescription()

			if test.resultErr != "" {
				if err == nil {
					t.Fatalf("getDescription() got nil error, expected %s", test.resultErr)
				}
				if err.Error() != test.resultErr {
					t.Fatalf("getDescription() got error %s, expected %s", err.Error(), test.resultErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("getDescription() got error %s, expected %s", err.Error(), test.result)
			}
			if desc != test.result {
				t.Fatalf("getDescription() got %s, expected %s", desc, test.result)
			}
		})
	}
}
