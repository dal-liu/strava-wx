package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DatabaseError struct {
	id types.AttributeValue
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Athlete id %v not found", e.id)
}

type DynamoDBClient struct {
	svc *dynamodb.Client
}

func (c DynamoDBClient) GetAccessToken(ctx context.Context, athleteId int) (AccessToken, error) {
	token := AccessToken{AthleteId: athleteId}
	err := c.getItem(ctx, token.GetKey(), "AccessTokens", &token)
	return token, err
}

func (c DynamoDBClient) UpdateAccessToken(ctx context.Context, token AccessToken) error {
	update := expression.Set(expression.Name("AccessToken"), expression.Value(token.Code))
	update.Set(expression.Name("ExpiresAt"), expression.Value(token.ExpiresAt))
	return c.updateItem(ctx, token.GetKey(), "AccessTokens", update)
}

func (c DynamoDBClient) GetRefreshToken(ctx context.Context, athleteId int) (RefreshToken, error) {
	token := RefreshToken{AthleteId: athleteId}
	err := c.getItem(ctx, token.GetKey(), "RefreshTokens", &token)
	return token, err
}

func (c DynamoDBClient) UpdateRefreshToken(ctx context.Context, token RefreshToken) error {
	update := expression.Set(expression.Name("RefreshToken"), expression.Value(token.Code))
	return c.updateItem(ctx, token.GetKey(), "RefreshTokens", update)
}

func (c DynamoDBClient) getItem(ctx context.Context, key map[string]types.AttributeValue, tableName string, out any) error {
	resp, err := c.svc.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	if len(resp.Item) == 0 {
		return &DatabaseError{key["AthleteId"]}
	}

	return attributevalue.UnmarshalMap(resp.Item, out)
}

func (c DynamoDBClient) updateItem(ctx context.Context, key map[string]types.AttributeValue, tableName string, update expression.UpdateBuilder) error {
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	_, err = c.svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	return err
}

func CreateClient(ctx context.Context) (DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return DynamoDBClient{}, err
	}
	return DynamoDBClient{dynamodb.NewFromConfig(cfg)}, nil
}
