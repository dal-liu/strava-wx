package strava

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type ActivityResponse struct {
	Start_date   string
	Start_latlng []float64
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
	return ar, json.NewDecoder(resp.Body).Decode(&ar)
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
	return err
}
