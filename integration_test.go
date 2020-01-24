package ddd

import (
	"fmt"
	"testing"
)

func _TestPublish(t *testing.T) {
	topic := "arn:aws:sns:us-west-2:202214144554:NotifyMe"
	sns,err := NewSnsClient(topic, "notification subject")
	if err != nil {
		panic(err)
	}
	out,err := sns.PublishMessage(`{"msg":"a test message in an object"}`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",out)
}

func _TestStorageService_Handle(t *testing.T) {
	//persist := NewDynamoClient("stc_interest_submission")
	//fmt.Println(err)
}