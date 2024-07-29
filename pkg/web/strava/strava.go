package strava

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type TokenResponse struct {
	Access_token  string
	Expires_at    int
	Refresh_token string
}

type ActivityResponse struct {
	Start_date   string
	Start_latlng []float64
}

func GetNewTokens(refreshToken string) (tr TokenResponse, err error) {
	req, err := http.NewRequest("POST", "https://www.strava.com/api/v3/oauth/token?grant_type=refresh_token", nil)
	if err != nil {
		return tr, err
	}

	q := req.URL.Query()
	q.Add("client_id", os.Getenv("STRAVA_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("STRAVA_CLIENT_SECRET"))
	q.Add("refresh_token", refreshToken)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tr, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return tr, err
	}

	return tr, nil
}

func GetActivity(activityId int, accessToken string) (ar ActivityResponse, err error) {
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/activities/"+strconv.Itoa(activityId), nil)
	if err != nil {
		return ar, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ar, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return ar, err
	}

	return ar, nil
}

func UpdateActivity(activityId int, accessToken string, description string) error {
	payload := []byte(`{"description":"` + description + `"}`)
	req, err := http.NewRequest("PUT", "https://www.strava.com/api/v3/activities/"+strconv.Itoa(activityId), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
