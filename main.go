package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

func handlePost(ctx context.Context, req events.LambdaFunctionURLRequest) (resp events.LambdaFunctionURLResponse, err error) {
	return resp, nil
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
