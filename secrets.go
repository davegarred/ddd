package ddd

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecret(id string) (map[string]interface{},error) {
	mgr, err := secretsClient("us-west-2")
	if err != nil {
		return nil,err
	}
	result, err := mgr.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     &id,
	})
	if err != nil {
		return nil,err
	}
	secrets := map[string]interface{}{}
	err = json.Unmarshal([]byte(*result.SecretString),&secrets)
	if err != nil {
		return nil,err
	}
	return secrets,nil
}

func secretsClient(region string) (*secretsmanager.SecretsManager, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return secretsmanager.New(sess), nil
}
