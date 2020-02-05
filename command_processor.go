package ddd

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

type CommandProcessor struct {
	supportedCommands map[string]func(context.Context,Request) events.APIGatewayProxyResponse
	*validator.Validate
	ut.Translator
}

func NewCommandProcessor() *CommandProcessor {
	processor := &CommandProcessor{
		supportedCommands: make(map[string]func(context.Context,Request) events.APIGatewayProxyResponse),
		Validate:          nil,
		Translator:        nil,
	}
	return processor
}

func (p *CommandProcessor) RegisterAll(supportedCommands map[string]func(context.Context, Request) events.APIGatewayProxyResponse) {
	p.supportedCommands = supportedCommands
}

func (p *CommandProcessor) Register(path string, f func(context.Context, Request) events.APIGatewayProxyResponse) {
	p.supportedCommands[path] = f
}

func (p *CommandProcessor) RegisterRpc(path string, f interface{}) error {
	wrapper, err := p.wrapRpcEndpoint(f)
	if err != nil {
		return err
	}
	p.supportedCommands[path] = func(ctx context.Context, req Request) events.APIGatewayProxyResponse {
		return wrapper(ctx, req.Body)
	}
	return nil
}

func (p *CommandProcessor) HandleRequest(ctx context.Context, reqEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := Request{ctx, reqEvent}
	commandFunction := p.supportedCommands[request.RequestPath(0)]
	if commandFunction == nil {
		fmt.Printf("bad route: %s\n", reqEvent.PathParameters)
		return ErrorResponse(errors.New("unsupported route")), nil
	}
	response := commandFunction(ctx, request)
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}
	return response, nil
}
