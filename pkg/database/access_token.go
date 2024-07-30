package database

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AccessToken struct {
	AthleteId int    `dynamodbav:"AthleteId"`
	Code      string `dynamodbav:"AccessToken"`
	ExpiresAt int    `dynamodbav:"ExpiresAt"`
}

func (a AccessToken) IsExpired() bool {
	return time.Now().Unix() >= int64(a.ExpiresAt)
}

func (a AccessToken) GetKey() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"AthleteId": &types.AttributeValueMemberN{Value: strconv.Itoa(a.AthleteId)},
	}
}
