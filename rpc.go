package ddd

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	en_locales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"reflect"
)

func (p *CommandProcessor) wrapRpcEndpoint(handlerFunc interface{}, v *validator.Validate) func(string) events.APIGatewayProxyResponse {
	handler := reflect.ValueOf(handlerFunc)
	validateHandler(handler)
	return p.buildHandlerWrapper(v, handler)
}

func (p *CommandProcessor) buildHandlerWrapper(v *validator.Validate, handler reflect.Value) func(body string) events.APIGatewayProxyResponse {
	dtoType := handler.Type().In(0)
	return func(body string) events.APIGatewayProxyResponse {
		dto := reflect.New(dtoType)
		if err := json.Unmarshal([]byte(body), dto.Interface()); err != nil {
			return ErrorResponse(err)
		}
		if v != nil {
			en := en_locales.New()
			translator := ut.New(en, en)
			englishTranslator, _ := translator.GetTranslator("en")
			en_translations.RegisterDefaultTranslations(v, englishTranslator)
			if err := v.Struct(dto.Interface()); err != nil {
				verrs,ok := err.(validator.ValidationErrors)
				if !ok {
				return ErrorResponse(err)
				}
				errors := make([]ErrorDetails,len(verrs))
				for i,fieldErr := range verrs {
					errors[i] = ErrorDetails{
						Field:   fieldErr.Field(),
						Tag:     fieldErr.Tag(),
						Message: fieldErr.Translate(englishTranslator),
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
	if  !handlerHasCorrectParameters {
		return errors.New("handler function interface is incorrect")
	}
	return nil
}
