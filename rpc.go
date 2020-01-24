package ddd

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

func (p *CommandProcessor) wrapRpcEndpoint(handlerFunc interface{}) func(string) events.APIGatewayProxyResponse {
	handler := reflect.ValueOf(handlerFunc)
	validateHandler(handler)
	return p.buildHandlerWrapper(handler)
}

func (p *CommandProcessor) buildHandlerWrapper(handler reflect.Value) func(body string) events.APIGatewayProxyResponse {
	dtoType := handler.Type().In(0)
	return func(body string) events.APIGatewayProxyResponse {
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

		inputValues := []reflect.Value{dto.Elem()}
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
	handlerHasCorrectParameters := handlerType.NumIn() == 1 && handlerType.NumOut() == 1
	if !handlerHasCorrectParameters {
		return errors.New("handler function interface is incorrect")
	}
	return nil
}
