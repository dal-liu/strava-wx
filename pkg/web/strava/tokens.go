package strava

import (
	"encoding/json"
	"net/http"
	"os"
)

type TokenResponse struct {
	Access_token  string
	Expires_at    int
	Refresh_token string
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
