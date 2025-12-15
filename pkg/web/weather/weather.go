package weather

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const epsilon float64 = 0.005
const mmPerInch float64 = 25.4

type weatherCondition struct {
	Id int
}

type weatherPrecipitation struct {
	One_hour float64 `json:"1h"`
}

type weatherData struct {
	Dt         int
	Sunrise    int
	Sunset     int
	Temp       float64
	Feels_like float64
	Humidity   int
	Wind_speed float64
	Wind_deg   int
	Wind_gust  float64
	Weather    []weatherCondition
	Rain       weatherPrecipitation
	Snow       weatherPrecipitation
}

func (wd weatherData) isDay() bool {
	return wd.Dt >= wd.Sunrise && wd.Dt < wd.Sunset
}

func (wd weatherData) getCondition() (string, error) {
	if len(wd.Weather) == 0 {
		return "", errors.New("No weather condition received")
	}

	switch wd.Weather[0].Id {
	case 200, 201, 202, 210, 211, 212, 221, 230, 231, 232:
		return "🌩️ Thunderstorm", nil
	case 300, 301, 302, 310, 311, 312, 313, 314, 321:
		return "🌧️ Drizzle", nil
	case 500, 501, 502, 503, 504, 511, 520, 521, 522, 531:
		return "🌧️ Rain", nil
	case 600, 601, 602, 611, 612, 613, 615, 616, 620, 621, 622:
		return "🌨️ Snow", nil
	case 701:
		return "🌫️ Mist", nil
	case 711:
		return "🌫️ Smoke", nil
	case 721:
		return "🌫️ Haze", nil
	case 731:
		return "🌫️ Dust", nil
	case 741:
		return "🌫️ Fog", nil
	case 751:
		return "🌫️ Sand", nil
	case 761:
		return "🌫️ Dust", nil
	case 762:
		return "🌫️ Ash", nil
	case 771:
		return "🌫️ Squall", nil
	case 781:
		return "🌪️ Tornado", nil
	case 800:
		if wd.isDay() {
			return "☀️ Sunny", nil
		} else {
			return "🌙 Clear", nil
		}
	case 801:
		if wd.isDay() {
			return "🌤️ Mostly sunny", nil
		} else {
			return "🌙 Mostly clear", nil
		}
	case 802:
		if wd.isDay() {
			return "⛅ Partly cloudy", nil
		} else {
			return "☁️ Partly cloudy", nil
		}
	case 803:
		if wd.isDay() {
			return "🌥️ Mostly cloudy", nil
		} else {
			return "☁️ Mostly cloudy", nil
		}
	case 804:
		return "☁️ Cloudy", nil
	}

	return "", errors.New("Unknown weather condition")
}

func (wd weatherData) getWindDirection() string {
	deg := float64(wd.Wind_deg)
	switch {
	case deg >= 348.75 || deg < 11.25:
		return "N"
	case deg >= 11.25 && deg < 33.75:
		return "NNE"
	case deg >= 33.75 && deg < 56.25:
		return "NE"
	case deg >= 56.25 && deg < 78.75:
		return "ENE"
	case deg >= 78.75 && deg < 101.25:
		return "E"
	case deg >= 101.25 && deg < 123.75:
		return "ESE"
	case deg >= 123.75 && deg < 146.25:
		return "SE"
	case deg >= 146.25 && deg < 168.75:
		return "SSE"
	case deg >= 168.75 && deg < 191.25:
		return "S"
	case deg >= 191.25 && deg < 213.75:
		return "SSW"
	case deg >= 213.75 && deg < 236.25:
		return "SW"
	case deg >= 236.25 && deg < 258.75:
		return "WSW"
	case deg >= 258.75 && deg < 281.25:
		return "W"
	case deg >= 281.25 && deg < 303.75:
		return "WNW"
	case deg >= 303.75 && deg < 326.25:
		return "NW"
	case deg >= 326.25 && deg < 348.75:
		return "NNW"
	}
	return "?"
}

func (wd weatherData) getPrecipitation() float64 {
	return (wd.Rain.One_hour + wd.Snow.One_hour) / mmPerInch
}

type weatherResponse struct {
	Data []weatherData
}

func (wr weatherResponse) getDescription() (string, error) {
	if len(wr.Data) == 0 {
		return "", errors.New("No weather data received")
	}
	data := wr.Data[0]

	var sb strings.Builder

	cond, err := data.getCondition()
	if err != nil {
		return "", err
	}
	sb.WriteString(cond)
	sb.WriteString(", ")

	sb.WriteString(strconv.FormatFloat(math.Round(data.Temp), 'f', -1, 64))
	sb.WriteString("°F, ")

	sb.WriteString("Feels like ")
	sb.WriteString(strconv.FormatFloat(math.Round(data.Feels_like), 'f', -1, 64))
	sb.WriteString("°F, ")

	sb.WriteString("Humidity ")
	sb.WriteString(strconv.Itoa(data.Humidity))
	sb.WriteString("%, ")

	sb.WriteString("Wind ")
	if windSpeed := math.Round(data.Wind_speed); windSpeed == 0 {
		sb.WriteString("0mph")
	} else {
		sb.WriteString(strconv.FormatFloat(windSpeed, 'f', -1, 64))
		sb.WriteString("mph ")

		if windGust := math.Round(data.Wind_gust); windGust > 0 {
			sb.WriteString("with ")
			sb.WriteString(strconv.FormatFloat(windGust, 'f', -1, 64))
			sb.WriteString("mph gusts ")
		}

		sb.WriteString("from ")
		sb.WriteString(data.getWindDirection())
	}

	if precip := data.getPrecipitation(); precip >= epsilon {
		sb.WriteString(", Precipitation ")
		sb.WriteString(strconv.FormatFloat(precip, 'f', 2, 64))
		sb.WriteString(" in/hr")
	}

	return sb.String(), nil
}

func GetWeatherDescription(lat, lon float64, date string) (string, error) {
	dt, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", "https://api.openweathermap.org/data/3.0/onecall/timemachine?units=imperial", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Add("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Add("dt", strconv.FormatInt(dt.Unix(), 10))
	q.Add("appid", os.Getenv("WEATHER_API_KEY"))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var wr weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&wr); err != nil {
		return "", err
	}

	return wr.getDescription()
}
