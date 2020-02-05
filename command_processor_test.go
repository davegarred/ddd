package ddd

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/onsi/gomega/types"
	"testing"

	. "github.com/onsi/gomega"
)

type testDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func testRpc(context.Context, testDto) events.APIGatewayProxyResponse {
	return Accepted()
}
func testRpc_missingContext(testDto) events.APIGatewayProxyResponse {
	return Accepted()
}
func testRpc_contextNotCorrectType(testDto, testDto) events.APIGatewayProxyResponse {
	return Accepted()
}
func testRpc_returnNotCorrectType(context.Context, testDto) testDto {
	return testDto{}
}

func TestCommandProcessor_RegisterRpc(t *testing.T) {
	handlerPath := "testFunction"

	tests := []struct {
		Name    string
		F       interface{}
		ErrorMatcher types.GomegaMatcher
		Matcher types.GomegaMatcher
	}{
		{
			Name:    "success",
			F:       testRpc,
			ErrorMatcher: BeNil(),
			Matcher: Not(BeNil()),
		},
		{
			Name: "missing context",
			F:    testRpc_missingContext,
			ErrorMatcher: Equal(errors.New("handler function interface is incorrect")),
			Matcher: BeNil(),
		},
		{
			Name: "first argument is not context.Context",
			F:    testRpc_contextNotCorrectType,
			ErrorMatcher: Equal(errors.New("handler function interface is incorrect - first argument must be of type context.Context")),
			Matcher: BeNil(),
		},
		{
			Name: "return argument is not events.APIGatewayProxyResponse",
			F:    testRpc_returnNotCorrectType,
			ErrorMatcher: Equal(errors.New("handler function interface is incorrect - return argument must be of type events.APIGatewayProxyResponse")),
			Matcher: BeNil(),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			comp := NewCommandProcessor()
			err := comp.RegisterRpc(handlerPath, test.F)
			NewGomegaWithT(t).Expect(comp.supportedCommands[handlerPath]).To(test.Matcher)
			NewGomegaWithT(t).Expect(err).To(test.ErrorMatcher)
		})
	}
}
