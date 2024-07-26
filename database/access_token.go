package database

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

func GetAccessToken(ctx context.Context, athleteId int) (AccessToken, error) {
	token := AccessToken{AthleteId: athleteId}
	err := getItem(ctx, token.GetKey(), "AccessTokens", &token)
	return token, err
}

func UpdateAccessToken(ctx context.Context, token AccessToken) error {
	update := expression.Set(expression.Name("AccessToken"), expression.Value(token.Code))
	update.Set(expression.Name("ExpiresAt"), expression.Value(token.ExpiresAt))
	return updateItem(ctx, token.GetKey(), "AccessTokens", update)
}
