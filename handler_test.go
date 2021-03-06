package ddd

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
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
