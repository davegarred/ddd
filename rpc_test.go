package ddd

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/onsi/gomega/types"
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

func TestBuildHandlerWrapper(t *testing.T) {
	h := &TestHandler{}
	commandProcessor := CommandProcessor{}
	err := commandProcessor.ConfigureValidator(map[string]string{
		"required":"{0} is required",
	})
	if err != nil {
		panic(err)
	}
	wrapper := commandProcessor.buildHandlerWrapper(reflect.ValueOf(h.SuccessfulTestMethod))
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
			expected: Equal(stdResponse(400, `{"error":"validation errors","validation_errors":[{"field":"Name","tag":"required","message":"Name is required"},{"field":"Birthdate","tag":"required","message":"Birthdate is required"}]}`)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := wrapper(context.Background(), test.body)
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

type TestHandler struct {
	Dto SampleDto
}

func (h *TestHandler) SuccessfulTestMethod(ctx context.Context, dto SampleDto) error {
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
