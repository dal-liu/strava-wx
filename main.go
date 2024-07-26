package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"strava-wx/database"
	"strava-wx/strava"
	"strava-wx/weather"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const athleteId int = 55440166

type webhookEvent struct {
	object_type string
	object_id   int
	aspect_type string
	owner_id    int
}

func handleGet(req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	if req.QueryStringParameters["hub.verify_token"] != os.Getenv("VERIFY_TOKEN") {
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	body, err := json.Marshal(map[string]string{"hub.challenge": req.QueryStringParameters["hub.challenge"]})
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	resp.StatusCode = http.StatusOK
	resp.Headers = map[string]string{"Content-Type": "application/json"}
	resp.Body = string(body)
	return resp, nil
}

func handlePost(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	go func() {
		if database.CreateClient(ctx) != nil {
			return
		}

		var event webhookEvent
		if json.Unmarshal([]byte(req.Body), &event) != nil {
			return
		}
		if event.owner_id != athleteId || event.object_type != "activity" || event.aspect_type != "create" {
			return
		}

		accessToken, err := database.GetAccessToken(ctx, athleteId)
		if err != nil {
			return
		}

		if accessToken.IsExpired() {
			getAndUpdateTokens(ctx, &accessToken)
		}
		activity, err := strava.GetActivity(event.object_id, accessToken.Code)
		if err != nil {
			return
		}

		description, err := weather.GetWeatherDescription(activity.Start_latlng[0], activity.Start_latlng[1], activity.Start_date)
		if err != nil {
			return
		}

		if accessToken.IsExpired() {
			getAndUpdateTokens(ctx, &accessToken)
		}
		strava.UpdateActivity(event.object_id, accessToken.Code, description)
	}()

	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK}, nil
}

func getAndUpdateTokens(ctx context.Context, accessToken *database.AccessToken) {
	refreshToken, err := database.GetRefreshToken(ctx, athleteId)
	newTokens, err := strava.GetNewTokens(refreshToken.Code)
	if err != nil {
		return
	}

	accessToken.Code = newTokens.Access_token
	accessToken.ExpiresAt = newTokens.Expires_at
	go database.UpdateAccessToken(ctx, *accessToken)

	refreshToken.Code = newTokens.Refresh_token
	go database.UpdateRefreshToken(ctx, refreshToken)
}

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if req.RequestContext.HTTP.Method == "GET" {
		return handleGet(req)
	}
	if req.RequestContext.HTTP.Method == "POST" {
		return handlePost(ctx, req)
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusMethodNotAllowed}, nil
}

func main() {
	lambda.Start(handler)
}
