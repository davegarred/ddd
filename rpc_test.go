package ddd

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/onsi/gomega/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"

	"testing"

	. "github.com/onsi/gomega"
)

var (
	expected = SampleDto{
		SampleId:  "id345z8w",
		Name:      "Joe Johnson",
		Birthdate: "2001-03-28",
	}
)

func TestValidateHandler(t *testing.T) {
	h := &TestHandler{}
	g := NewGomegaWithT(t)

	tests := []struct {
		name     string
		method   interface{}
		expected types.GomegaMatcher
	}{
		{
			name:     "success",
			method:   h.SuccessfulTestMethod,
			expected: BeNil(),
		},
		{
			name:     "invalid - no error",
			method:   InvalidFunc_NoDto,
			expected: Equal(errors.New("handler function interface is incorrect")),
		},
		{
			name:     "invalid - not dto processor",
			method:   InvalidFunc_NoError,
			expected: Equal(errors.New("handler function interface is incorrect")),
		},
		{
			name:     "invalid - not a function",
			method:   SampleDto{},
			expected: Equal(errors.New("handler function interface is incorrect")),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := reflect.ValueOf(test.method)
			err := validateHandler(handler)
			g.Expect(err).To(test.expected)
		})
	}
}

func InvalidFunc_NoError(dto SampleDto) {}
func InvalidFunc_NoDto() error          { return nil }

func TestBuildHandlerWrapper(t *testing.T) {
	h := &TestHandler{}
	commandProcessor := CommandProcessor{}
	wrapper := commandProcessor.buildHandlerWrapper(validator.New(), reflect.ValueOf(h.SuccessfulTestMethod))
	g := NewGomegaWithT(t)

	tests := []struct {
		name     string
		body     string
		expected types.GomegaMatcher
	}{
		{
			name:     "happy path",
			body:     `{"id":"id345z8w","name":"Joe Johnson","birthdate":"2001-03-28"}`,
			expected: Equal(stdResponse(202, ``)),
		},
		{
			name:     "bad json",
			body:     ``,
			expected: Equal(stdResponse(400, `{"error":"unexpected end of JSON input"}`)),
		},
		{
			name:     "validation errors",
			body:     `{"id":"id345z8w","name":"","birthdate":""}`,
			expected: Equal(stdResponse(400, `{"error":"validation errors","validation_errors":[{"field":"Name","tag":"required","message":"Name is a required field"},{"field":"Birthdate","tag":"required","message":"Birthdate is a required field"}]}`)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := wrapper(test.body)
			g.Expect(result).To(test.expected)
		})
	}
}

func stdResponse(status int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
	}
}

//func TestHandleBody_validationError(t *testing.T) {
//	h := &TestHandler{}
//
//	err := HandleBody(h.SuccessfulTestMethod, validator.New())([]byte(`{"id":"id345z8w"}`))
//
//	g := NewGomegaWithT(t)
//	errs := err.(validator.ValidationErrors)
//	g.Expect(len(errs)).To(Equal(2))
//	g.Expect(h.Dto).To(Equal(SampleDto{}))
//}

type TestHandler struct {
	Dto SampleDto
}

func (h *TestHandler) SuccessfulTestMethod(dto SampleDto) error {
	h.Dto = dto
	return nil
}

func (h *TestHandler) ErrorTestMethod(dto SampleDto) error {
	return errors.New("some error")
}

type SampleId string
type SampleDto struct {
	SampleId  `json:"id"`
	Name      string `json:"name" validate:"required"`
	Birthdate string `json:"birthdate" validate:"required"`
}
