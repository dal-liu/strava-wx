package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"strava-wx/pkg/database"
	"strava-wx/pkg/web/strava"
	"strava-wx/pkg/web/weather"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const athleteId int = 55440166

type webhookEvent struct {
	Object_type string
	Object_id   int
	Aspect_type string
	Owner_id    int
	Event_time  int
}

func workerHandler(ctx context.Context, req events.SQSEvent) error {
	log.Println("Received POST request. Creating DynamoDB client...")
	if err := database.CreateClient(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(req.Records))
	wg.Add(len(req.Records))
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	log.Println("Client created. Processing messages...")
	for i, record := range req.Records {
		log.Printf("Processing record %d...", i)
		go func(i int, record events.SQSMessage) {
			if err := processRecord(ctx, record); err != nil {
				errorChan <- err
			} else {
				log.Printf("Record %d processed.", i)
			}
			wg.Done()
		}(i, record)
	}

	for err := range errorChan {
		log.Println("ERROR:", err)
		return err
	}

	log.Println("All messages processed.")
	return nil
}

func processRecord(ctx context.Context, record events.SQSMessage) error {
	log.Println("Parsing record...")
	var event webhookEvent
	if err := json.Unmarshal([]byte(record.Body), &event); err != nil {
		log.Println("ERROR:", err)
		return err
	}

	log.Println("Record parsed. Checking if event is activity creation...")
	if event.Owner_id == athleteId && event.Object_type == "activity" && event.Aspect_type == "create" {
		log.Println("Event is activity creation. Getting access token...")
		accessToken, err := database.GetAccessToken(ctx, athleteId)
		if err != nil {
			log.Println("ERROR:", err)
			return err
		}

		log.Println("Retrieved access token.")
		if err = checkAccessToken(ctx, &accessToken); err != nil {
			log.Println("ERROR:", err)
			return err
		}
		log.Println("Getting activity...")
		activity, err := strava.GetActivity(event.Object_id, accessToken.Code)
		if err != nil {
			log.Println("ERROR:", err)
			return err
		}

		log.Println("Activity retrieved. Checking if activity has start coordinates...")
		if len(activity.Start_latlng) == 2 {
			log.Println("Activity has start coordinates. Getting weather description...")
			description, err := weather.GetWeatherDescription(activity.Start_latlng[0], activity.Start_latlng[1], activity.Start_date)
			if err != nil {
				log.Println("ERROR:", err)
				return err
			}

			log.Println("Weather description retrieved. Updating activity...")
			if err = checkAccessToken(ctx, &accessToken); err != nil {
				log.Println("ERROR:", err)
				return err
			}
			log.Println("Updating activity...")
			if err = strava.UpdateActivity(event.Object_id, accessToken.Code, description); err != nil {
				log.Println("ERROR:", err)
				return err
			}
			log.Println("Activity updated.")
		} else {
			log.Println("Activity does not have start coordinates. Ignoring...")
		}
	} else {
		log.Println("Event is not activity creation. Ignoring...")
	}

	return nil
}

func checkAccessToken(ctx context.Context, accessToken *database.AccessToken) error {
	log.Println("Checking if access token is expired...")
	if accessToken.IsExpired() {
		log.Println("Access token is expired. Getting refresh token...")
		refreshToken, err := database.GetRefreshToken(ctx, athleteId)
		if err != nil {
			log.Println("ERROR:", err)
			return err
		}
		log.Println("Refresh token retrieved. Getting new tokens...")
		newTokens, err := strava.GetNewTokens(refreshToken.Code)
		if err != nil {
			log.Println("ERROR:", err)
			return err
		}

		var wg sync.WaitGroup
		errorChan := make(chan error, 2)
		wg.Add(2)
		go func() {
			wg.Wait()
			close(errorChan)
		}()

		log.Println("New tokens retrieved. Updating tokens...")

		accessToken.Code = newTokens.Access_token
		accessToken.ExpiresAt = newTokens.Expires_at
		go func() {
			if err = database.UpdateAccessToken(ctx, *accessToken); err != nil {
				errorChan <- err
			} else {
				log.Println("Access token updated.")
			}
			wg.Done()
		}()

		refreshToken.Code = newTokens.Refresh_token
		go func() {
			if err = database.UpdateRefreshToken(ctx, refreshToken); err != nil {
				errorChan <- err
			} else {
				log.Println("Refresh token updated.")
			}
			wg.Done()
		}()

		for err := range errorChan {
			log.Println("ERROR:", err)
			return err
		}

		log.Println("Tokens updated.")
	}
	return nil
}

func main() {
	lambda.Start(workerHandler)
}
