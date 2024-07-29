package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"strava-wx/pkg/queue"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleGet(req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	log.Println("Received GET request. Verifying token...")
	if req.QueryStringParameters["hub.verify_token"] != os.Getenv("VERIFY_TOKEN") {
		log.Println("Token verification failed.")
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	log.Println("Token verified. Extracting challenge...")
	body, err := json.Marshal(map[string]string{"hub.challenge": req.QueryStringParameters["hub.challenge"]})
	if err != nil {
		log.Println("ERROR:", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Println("Challenge extracted. Responding OK...")
	resp.StatusCode = http.StatusOK
	resp.Headers = map[string]string{"Content-Type": "application/json"}
	resp.Body = string(body)
	return resp, nil
}

func handlePost(ctx context.Context, req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	log.Println("Received POST request. Creating SQS client...")
	if err = queue.CreateClient(ctx); err != nil {
		log.Println("ERROR:", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Println("Client created. Sending message to queue...")
	if err = queue.Send(ctx, req.Body, os.Getenv("QUEUE_URL")); err != nil {
		log.Println("ERROR:", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Println("Message sent. Responding OK...")
	resp.StatusCode = http.StatusOK
	return resp, nil
}

func webhookHandler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if req.RequestContext.HTTP.Method == "GET" {
		return handleGet(req)
	}
	if req.RequestContext.HTTP.Method == "POST" {
		return handlePost(ctx, req)
	}
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusMethodNotAllowed}, nil
}

func main() {
	lambda.Start(webhookHandler)
}
