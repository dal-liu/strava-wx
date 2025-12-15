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
			input:     weatherData{Weather: []weatherCondition{{Id: 0}}},
			resultErr: "Unknown weather condition",
		},
		"cloudy": {
			input:  weatherData{Weather: []weatherCondition{{Id: 804}}},
			result: "☁️ Cloudy",
		},
		"sunny": {
			input:  weatherData{Dt: 2, Sunrise: 1, Sunset: 3, Weather: []weatherCondition{{Id: 800}}},
			result: "☀️ Sunny",
		},
		"clear": {
			input:  weatherData{Dt: 4, Sunrise: 1, Sunset: 3, Weather: []weatherCondition{{Id: 800}}},
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
			input:  weatherData{Rain: weatherPrecipitation{One_hour: 12.7}, Snow: weatherPrecipitation{One_hour: 12.7}},
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
			input:     weatherResponse{Data: []weatherData{{}}},
			resultErr: "No weather condition received",
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
