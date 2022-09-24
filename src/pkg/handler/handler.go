package handler

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	AccessTokenURL string
	SubredditURL   string
}

func NewHandler(accessTokenURL string, subredditURL string) (*Handler, error) {
	h := Handler{
		AccessTokenURL: accessTokenURL,
		SubredditURL:   subredditURL,
	}
	return &h, nil
}

func (h *Handler) Handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return h.response("", http.StatusOK), nil
}

func (h *Handler) response(body string, statusCode int) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers:    map[string]string{},
	}
	return resp
}
