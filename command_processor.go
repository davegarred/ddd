package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type CommandProcessor struct {
	supportedCommands map[string]func(context.Context, events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

func NewCommandProcessor(supportedCommands map[string]func(context.Context, events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) *CommandProcessor {
	processor := &CommandProcessor{supportedCommands}
	return processor
}

func (p *CommandProcessor) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lContext, ok := lambdacontext.FromContext(ctx)
	if !ok {
		errMsg := fmt.Sprintf("invalid context encountered: %+v", ctx)
		err := errors.New(errMsg)
		return ErrorResponse(err), err
	}
	fmt.Printf("cognito identity: %v\n", lContext.Identity)

	route := request.PathParameters["proxy"]
	commandFunction := p.supportedCommands[route]
	if commandFunction == nil {
		fmt.Printf("bad route: %s\n", route)
		return ErrorResponse(errors.New("unsupported route")), nil
	}

	if ser, err := json.Marshal(request); err == nil {
		fmt.Println("request body")
		fmt.Println(string(ser))
	} else {
		fmt.Printf("error - %v\n", err)
	}

	if request.Headers["Content-Type"] == "application/json" {
		fmt.Println(request.Body)
	}

	response := commandFunction(ctx, request)
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}
	return response, nil
}
