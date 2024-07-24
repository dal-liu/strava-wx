package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

const athleteId int = 55440166

type TokenResponse struct {
	Access_token  string
	Expires_at    int
	Refresh_token string
}

func refreshExpiredTokens(refreshToken string) (TokenResponse, error) {
	req, err := http.NewRequest("POST", "https://www.strava.com/oauth/token?grant_type=refresh_token", nil)
	if err != nil {
		return TokenResponse{}, err
	}

	q := req.URL.Query()
	q.Add("client_id", os.Getenv("STRAVA_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("STRAVA_CLIENT_SECRET"))
	q.Add("refresh_token", refreshToken)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}

	defer resp.Body.Close()

	var body TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return TokenResponse{}, err
	}

	return body, nil
}

type ActivityResponse struct {
	Start_date   string
	Start_latlng []float64
}

func updateActivity(id int) error {
	req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/activities/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	token, err := getAccessToken(athleteId)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var body ActivityResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	lat, lon := body.Start_latlng[0], body.Start_latlng[1]
	date := body.Start_date
	desc, err := getDescription(lat, lon, date)
	if err != nil {
		return err
	}

	payload := []byte(`{"description":"` + desc + `"}`)
	req, err = http.NewRequest("PUT", "https://www.strava.com/api/v3/activities/"+strconv.Itoa(id), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	token, err = getAccessToken(athleteId)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
