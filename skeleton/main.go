package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davegarred/ddd"
)

func main() {
	processor := ddd.NewCommandProcessor()
	processor.Register("target", Handler)
	lambda.Start(processor.HandleRequest)
}

func Handler(_ context.Context, req ddd.Request) events.APIGatewayProxyResponse {
	//if ser, err := json.Marshal(req.APIGatewayProxyRequest); err != nil {
	//	fmt.Println(string(ser))
	//}
	return ddd.Ok(req)
}
