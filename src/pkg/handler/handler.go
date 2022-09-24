package handler

import (
	"net/http"
	"wr-latest-daily-redirect-serverless/pkg/daily"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	AccessTokenURL string
	ClientId       string
	ClientSecret   string
	SubredditURL   string
}

func NewHandler(clientId string, clientSecret string, accessTokenURL string, subredditURL string) (*Handler, error) {
	h := Handler{
		AccessTokenURL: accessTokenURL,
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		SubredditURL:   subredditURL,
	}
	return &h, nil
}

func (h *Handler) Handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	d := daily.NewDaily(h.ClientId, h.ClientSecret, "golang:wr-latest-daily-redirect-serverless:1.0.0 (by /u/murrtu)", h.AccessTokenURL, h.SubredditURL)

	posts, err := d.Posts()
	if err != nil {
		return h.response(err.Error(), 500, nil), err
	}

	// Latest
	latest, err := d.Latest(posts)
	if err != nil {
		return h.response(err.Error(), 500, nil), err
	}

	return h.response(latest.Url, http.StatusTemporaryRedirect, &latest.Url), nil
}

func (h *Handler) response(body string, statusCode int, redirectURL *string) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers:    map[string]string{},
	}
	if redirectURL != nil {
		resp.Headers["Location"] = *redirectURL
	}
	return resp
}
