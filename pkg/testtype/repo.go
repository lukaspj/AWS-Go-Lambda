package testtype

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoRepo struct {
	TableName string `json:"table_name"`
	client    *dynamodb.Client
}

func NewDynamoRepo(config aws.Config, tableName string) *DynamoRepo {
	return &DynamoRepo{
		TableName: tableName,
		client:    dynamodb.NewFromConfig(config),
	}
}

type TestItem struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

func (d *DynamoRepo) Create(ctx context.Context, t *TestItem) (string, error) {
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"UserId": &types.AttributeValueMemberS{Value: uuid.New().String()},
			"Name":   &types.AttributeValueMemberS{Value: t.Name},
		},
		TableName: aws.String(d.TableName),
	})
	return uuid.New().String(), err
}

func (d *DynamoRepo) Get(ctx context.Context, id string) (*TestItem, error) {
	item, err := d.client.GetItem(ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String("test_table"),
			Key: map[string]types.AttributeValue{
				"UserId": &types.AttributeValueMemberS{Value: id},
			},
		})
	if err != nil {
		return nil, err
	}

	return &TestItem{
		UserId: item.Item["UserId"].(*types.AttributeValueMemberS).Value,
		Name:   item.Item["Name"].(*types.AttributeValueMemberS).Value,
	}, nil
}

func (d *DynamoRepo) List(ctx context.Context) ([]TestItem, error) {
	scan, err := d.client.Scan(ctx,
		&dynamodb.ScanInput{
			TableName: aws.String("test_table"),
		})
	if err != nil {
		return nil, err
	}

	items := make([]TestItem, len(scan.Items))

	for i, item := range scan.Items {
		items[i] = TestItem{
			UserId: item["UserId"].(*types.AttributeValueMemberS).Value,
			Name:   item["Name"].(*types.AttributeValueMemberS).Value,
		}
	}

	return items, nil
}
