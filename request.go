package ddd

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type Request struct {
	context.Context
	events.APIGatewayProxyRequest
}

func (r *Request) RequestPath(section int) string {
	route := strings.Split(r.PathParameters["proxy"], "/")
	if section > len(route) {
		return ""
	}
	return route[section]
}

func (r *Request) Claims() map[string]interface{} {
	return r.RequestContext.Authorizer["claims"].(map[string]interface{})
}

func (r *Request) CognitoUsername() string {
	claims := r.Claims()
	return claims["cognito:username"].(string)
}
func (r *Request) CognitoGroups() string {
	claims := r.Claims()
	groups := claims["cognito:groups"]
	if groups != nil {
		return groups.(string)
	}
	return ""
}
