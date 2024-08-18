package weather

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type weatherResponse struct {
	Data []struct {
		Dt         int
		Sunrise    int
		Sunset     int
		Temp       float64
		Feels_like float64
		Humidity   int
		Wind_speed float64
		Wind_deg   int
		Wind_gust  float64
		Weather    []struct {
			Id int
		}
		Rain struct {
			One_hour float64 `json:"1h"`
		}
		Snow struct {
			One_hour float64 `json:"1h"`
		}
	}
}

func (wr weatherResponse) getDescription() string {
	if len(wr.Data) == 0 {
		return ""
	}

	var sb strings.Builder
	data := wr.Data[0]

	sb.WriteString(wr.getCondition())
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
		sb.WriteString(wr.getWindDirection())
	}

	if precip := wr.getPrecipitation(); precip >= 0.005 {
		sb.WriteString(", Precipitation ")
		sb.WriteString(strconv.FormatFloat(precip, 'f', 2, 64))
		sb.WriteString(" in/hr")
	}

	return sb.String()
}

func (wr weatherResponse) getCondition() string {
	if len(wr.Data) > 0 && len(wr.Data[0].Weather) > 0 {
		switch wr.Data[0].Weather[0].Id {
		case 200, 201, 202, 210, 211, 212, 221, 230, 231, 232:
			return "🌩️ Thunderstorm"
		case 300, 301, 302, 310, 311, 312, 313, 314, 321:
			return "🌧️ Drizzle"
		case 500, 501, 502, 503, 504, 511, 520, 521, 522, 531:
			return "🌧️ Rain"
		case 600, 601, 602, 611, 612, 613, 615, 616, 620, 621, 622:
			return "🌨️ Snow"
		case 701:
			return "🌫️ Mist"
		case 711:
			return "🌫️ Smoke"
		case 721:
			return "🌫️ Haze"
		case 731:
			return "🌫️ Dust"
		case 741:
			return "🌫️ Fog"
		case 751:
			return "🌫️ Sand"
		case 761:
			return "🌫️ Dust"
		case 762:
			return "🌫️ Ash"
		case 771:
			return "🌫️ Squall"
		case 781:
			return "🌪️ Tornado"
		case 800:
			if wr.isDay() {
				return "☀️ Sunny"
			} else {
				return "🌙 Clear"
			}
		case 801:
			if wr.isDay() {
				return "🌤️ Mostly sunny"
			} else {
				return "🌙 Mostly clear"
			}
		case 802:
			if wr.isDay() {
				return "⛅ Partly cloudy"
			} else {
				return "☁️ Partly cloudy"
			}
		case 803:
			if wr.isDay() {
				return "🌥️ Mostly cloudy"
			} else {
				return "☁️ Mostly cloudy"
			}
		case 804:
			return "☁️ Cloudy"
		}
	}
	return ""
}

func (wr weatherResponse) isDay() bool {
	if len(wr.Data) == 0 {
		return false
	}
	dt := wr.Data[0].Dt
	sunrise := wr.Data[0].Sunrise
	sunset := wr.Data[0].Sunset
	return dt >= sunrise && dt < sunset
}

func (wr weatherResponse) getWindDirection() string {
	if len(wr.Data) > 0 {
		deg := float64(wr.Data[0].Wind_deg)
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
	}
	return ""
}

func (wr weatherResponse) getPrecipitation() float64 {
	if len(wr.Data) == 0 {
		return 0
	}
	return (wr.Data[0].Rain.One_hour + wr.Data[0].Snow.One_hour) / 25.4
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

	return wr.getDescription(), nil
}
