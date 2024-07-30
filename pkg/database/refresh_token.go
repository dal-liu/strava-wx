package database

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type RefreshToken struct {
	AthleteId int    `dynamodbav:"AthleteId"`
	Code      string `dynamodbav:"RefreshToken"`
}

func (r RefreshToken) GetKey() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"AthleteId": &types.AttributeValueMemberN{Value: strconv.Itoa(r.AthleteId)},
	}
}
