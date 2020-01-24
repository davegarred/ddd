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
	supportedCommands map[string]func(Request) events.APIGatewayProxyResponse
	*validator.Validate
	ut.Translator
}

func NewCommandProcessor() *CommandProcessor {
	processor := &CommandProcessor{
		supportedCommands: make(map[string]func(request Request) events.APIGatewayProxyResponse),
		Validate:          nil,
		Translator:        nil,
	}
	return processor
}

func (p *CommandProcessor) RegisterAll(supportedCommands map[string]func(Request) events.APIGatewayProxyResponse) {
	p.supportedCommands = supportedCommands
}

func (p *CommandProcessor) Register(path string, f func(Request) events.APIGatewayProxyResponse) {
	p.supportedCommands[path] = f
}

func (p *CommandProcessor) RegisterRpc(path string, f interface{}) {
	wrapper := p.wrapRpcEndpoint(f)
	p.supportedCommands[path] = func(req Request) events.APIGatewayProxyResponse {
		return wrapper(req.Body)
	}
}

func (p *CommandProcessor) HandleRequest(ctx context.Context, reqEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := Request{ctx, reqEvent}
	commandFunction := p.supportedCommands[request.RequestPath(0)]
	response := commandFunction(request)
	if commandFunction == nil {
		fmt.Printf("bad route: %s\n", reqEvent.PathParameters)
		return ErrorResponse(errors.New("unsupported route")), nil
	}
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}
	return response, nil
}
