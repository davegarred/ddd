package ddd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

func (p *CommandProcessor) wrapRpcEndpoint(handlerFunc interface{}) (func(context.Context, string) events.APIGatewayProxyResponse, error) {
	handler := reflect.ValueOf(handlerFunc)
	if err := validateHandler(handler); err != nil {
		return nil, err
	}
	return p.buildHandlerWrapper(handler), nil
}

func (p *CommandProcessor) buildHandlerWrapper(handler reflect.Value) func(ctx context.Context, body string) events.APIGatewayProxyResponse {
	dtoType := handler.Type().In(1)
	return func(ctx context.Context, body string) events.APIGatewayProxyResponse {
		contextValue := reflect.ValueOf(ctx)
		dto := reflect.New(dtoType)
		if err := json.Unmarshal([]byte(body), dto.Interface()); err != nil {
			return ErrorResponse(err)
		}
		if p.Validate != nil {
			if err := p.Struct(dto.Interface()); err != nil {
				verrs, ok := err.(validator.ValidationErrors)
				if !ok {
					return ErrorResponse(err)
				}
				errors := make([]ErrorDetails, len(verrs))
				for i, fieldErr := range verrs {
					errors[i] = ErrorDetails{
						Field:   fieldErr.Field(),
						Tag:     fieldErr.Tag(),
						Message: fieldErr.Translate(p.Translator),
					}
				}
				return ValidationErrorResponse(&errors)
			}
		}
		inputValues := []reflect.Value{contextValue, dto.Elem()}
		outputValues := handler.Call(inputValues)
		result := outputValues[0].Interface()
		if result == nil {
			return Accepted()
		}
		return result.(events.APIGatewayProxyResponse)
	}
}

func validateHandler(handler reflect.Value) error {
	handlerType := handler.Type()
	if handler.Kind() != reflect.Func {
		return errors.New("handler function interface is incorrect")
	}
	handlerHasCorrectParameters := handlerType.NumIn() == 2 && handlerType.NumOut() == 1
	if !handlerHasCorrectParameters {
		return errors.New("handler function interface is incorrect")
	}
	contextType := handler.Type().In(0)
	if _, ok := reflect.New(contextType).Interface().(*context.Context); !ok {
		return errors.New("handler function interface is incorrect - first argument must be of type context.Context")
	}
	returnType := handler.Type().Out(0)
	if _, ok := reflect.New(returnType).Interface().(*events.APIGatewayProxyResponse); !ok {
		return errors.New("handler function interface is incorrect - return argument must be of type events.APIGatewayProxyResponse")
	}
	return nil
}
