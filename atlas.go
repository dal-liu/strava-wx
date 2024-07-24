package main

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenDocument struct {
	Athlete_id    int
	Access_token  string
	Expires_at    int
	Refresh_token string
}

var client *mongo.Client

func connectToMongoDB() error {
	var err error

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("MONGODB_URI is not set")
	}
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	return nil
}

func disconnectFromMongoDB() error {
	if err := client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}

func getAccessToken(athleteId int) (string, error) {
	accessColl := client.Database("authDB").Collection("access_tokens")
	filter := bson.M{"athlete_id": athleteId}
	var doc TokenDocument
	if err := accessColl.FindOne(context.TODO(), filter).Decode(&doc); err != nil {
		return "", err
	}

	var token string

	if doc.Expires_at <= int(time.Now().Unix()) {
		resp, err := refreshExpiredTokens(doc.Refresh_token)
		if err != nil {
			return "", err
		}

		accessUpdate := bson.M{"$set": bson.M{"access_token": resp.Access_token, "expires_at": resp.Expires_at}}
		go accessColl.UpdateOne(context.TODO(), filter, accessUpdate)

		refreshUpdate := bson.M{"$set": bson.M{"refresh_token": resp.Refresh_token}}
		go func() {
			refreshColl := client.Database("authDB").Collection("refresh_tokens")
			refreshColl.UpdateOne(context.TODO(), filter, refreshUpdate)
		}()

		token = resp.Access_token
	} else {
		token = doc.Access_token
	}

	return token, nil
}
