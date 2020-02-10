package ddd

import (
	"github.com/aws/aws-lambda-go/events"
)

func Ok(payload interface{}) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       cleanlySerialize(payload),
	}
}
func Accepted() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 202,
	}
}
func NotAuthorized() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 401,
	}
}
func NotFound() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
	}
}
func ErrorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       cleanlySerialize(ErrorDto{err.Error(), nil}),
	}
}
func ValidationErrorResponse(err *[]ErrorDetails) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       cleanlySerialize(ErrorDto{"validation errors", err}),
	}
}
func ServerErrorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       cleanlySerialize(ErrorDto{err.Error(), nil}),
	}
}
