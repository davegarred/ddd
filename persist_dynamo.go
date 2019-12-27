package ddd

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	primaryKey   = "id"
	payloadField = "payload"
)

var _ Persist = (*DynamoClient)(nil)

func (c *DynamoClient) Load(id string, object interface{}) error {
	result, err := c.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			primaryKey: {
				S: aws.String(id),
			},
		},
		TableName: aws.String(c.tableName),
	})
	if err != nil {
		return err
	}
	payload := result.Item[payloadField]
	if payload == nil {
		return errors.New("not found")
	}
	payloadString := payload.S
	return json.Unmarshal([]byte(*payloadString), object)
}

func (c *DynamoClient) Save(id string, payload interface{}) error {
	ser,err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			primaryKey: {
				S: aws.String(id),
			},
			payloadField: {
				S: aws.String(string(ser)),
			}},
		TableName: aws.String(c.tableName),
	})
	return err
}

type DynamoClient struct {
	svc       *dynamodb.DynamoDB
	tableName string
}

func NewDynamoClient(tableName string) *DynamoClient {
	table, err := client("us-west-2")
	if err != nil {
		panic(err)
	}
	return &DynamoClient{
		svc:                table,
		tableName: tableName,
	}
}

func client(region string) (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}
