package database

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

func GetRefreshToken(ctx context.Context, athleteId int) (RefreshToken, error) {
	token := RefreshToken{AthleteId: athleteId}
	err := getItem(ctx, token.GetKey(), "RefreshTokens", &token)
	return token, err
}

func UpdateRefreshToken(ctx context.Context, token RefreshToken) error {
	update := expression.Set(expression.Name("RefreshToken"), expression.Value(token.Code))
	return updateItem(ctx, token.GetKey(), "RefreshTokens", update)
}
