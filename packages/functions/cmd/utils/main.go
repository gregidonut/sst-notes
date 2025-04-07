package utils

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func GetUserId(event events.APIGatewayProxyRequest) (string, error) {

	type Jwt struct {
		Claims map[string]interface{} `json:"claims"`
	}
	jj := Jwt{}

	if j, ok := event.RequestContext.Authorizer["jwt"]; ok {
		b, err := json.Marshal(j)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(b, &jj); err != nil {
			return "", err
		}
	}

	userId := jj.Claims["sub"]
	return userId.(string), nil
}
