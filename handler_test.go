package ddd

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/onsi/gomega/types"
	"testing"

	. "github.com/onsi/gomega"
)

func test(context.Context, Request) events.APIGatewayProxyResponse {
	return Accepted()
}

func TestHandlers(t *testing.T) {
	comp := NewCommandProcessor()
	comp.Register("testFunction", test)

	tests := []struct {
		Name string
		lambdacontext.LambdaContext
		events.APIGatewayProxyRequest
		ExpectedResponse int
	}{
		{
			Name:          "testFunction",
			LambdaContext: lambdacontext.LambdaContext{},
			APIGatewayProxyRequest: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"proxy": "testFunction"},
				Body:           "{}",
			},
			ExpectedResponse: 202,
		},
		{
			Name:                   "submit - no data",
			LambdaContext:          lambdacontext.LambdaContext{},
			APIGatewayProxyRequest: events.APIGatewayProxyRequest{},
			ExpectedResponse:       400,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctx := lambdacontext.NewContext(context.Background(), &test.LambdaContext)
			resp, _ := comp.HandleRequest(ctx, test.APIGatewayProxyRequest)

			NewGomegaWithT(t).Expect(resp.StatusCode).To(Equal(test.ExpectedResponse))
		})
	}
}

type testDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func testRpc(context.Context, testDto) events.APIGatewayProxyResponse {
	return Accepted()
}
func testRpc_missingContext(context.Context, testDto) events.APIGatewayProxyResponse {
	return Accepted()
}

func TestCommandProcessor_RegisterRpc(t *testing.T) {
	comp := NewCommandProcessor()
	handlerPath := "testFunction"

	tests := []struct {
		Name    string
		F       interface{}
		Matcher types.GomegaMatcher
	}{
		{
			Name:    "success",
			F:       testRpc,
			Matcher: Not(BeNil()),
		},
		{
			Name: "missing context",
			F:    testRpc_missingContext,
			Matcher: BeNil(),
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			comp.RegisterRpc(handlerPath, test.F)
			NewGomegaWithT(t).Expect(comp.supportedCommands[handlerPath]).To(test.Matcher)
		})
	}
}
