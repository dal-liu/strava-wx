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

var client *mongo.Client

type AccessDocument struct {
	Access_token string
	Expires_at   int
}

type RefreshDocument struct {
	Refresh_token string
}

func (td AccessDocument) expired() bool {
	return td.Expires_at <= int(time.Now().Unix())
}

func connectToMongoDB() (err error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("MONGODB_URI is not set")
	}
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return
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
	filter := bson.M{"athlete_id": athleteId}

	accessColl := client.Database("authDB").Collection("access_tokens")
	var ad AccessDocument
	if err := accessColl.FindOne(context.TODO(), filter).Decode(&ad); err != nil {
		return "", err
	}

	var returnToken string

	if ad.expired() {
		refreshColl := client.Database("authDB").Collection("refresh_tokens")
		var rd RefreshDocument
		if err := refreshColl.FindOne(context.TODO(), filter).Decode(&rd); err != nil {
			return "", err
		}

		accessToken, expiresAt, refreshToken, err := refreshExpiredTokens(rd.Refresh_token)
		if err != nil {
			return "", err
		}

		accessUpdate := bson.M{"$set": bson.M{"access_token": accessToken, "expires_at": expiresAt}}
		go accessColl.UpdateOne(context.TODO(), filter, accessUpdate)

		refreshUpdate := bson.M{"$set": bson.M{"refresh_token": refreshToken}}
		if rd.Refresh_token != refreshToken {
			go refreshColl.UpdateOne(context.TODO(), filter, refreshUpdate)
		}

		returnToken = accessToken
	} else {
		returnToken = ad.Access_token
	}

	return returnToken, nil
}
