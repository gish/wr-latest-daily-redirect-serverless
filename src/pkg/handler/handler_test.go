package handler

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleRespondsWithFailWhenGetAccessTokenFails(t *testing.T) {
	req := events.APIGatewayProxyRequest{}
	expected := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	resp, err := Handle(req)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
	}
	if resp.StatusCode != expected.StatusCode {
		t.Errorf("request did not fail, status code %d", resp.StatusCode)
	}

}

func TestHandleRespondsWithFailWhenGetPostsFail(t *testing.T) {
}

func TestHandleRespondsWithCodeAndLocationWhenOK(t *testing.T) {
}
