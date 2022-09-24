package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleRespondsWithFailWhenGetAccessTokenFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	h, err := NewHandler("clientId", "clientSecret", server.URL, "")
	if err != nil {
		t.Errorf("failed with error: %+v", err)
	}

	req := events.APIGatewayProxyRequest{}
	expected := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	resp, err := h.Handle(req)
	if err == nil {
		t.Errorf("got no expected error")
	}
	if resp.StatusCode != expected.StatusCode {
		t.Errorf("request did not fail, status code %d", resp.StatusCode)
	}
}

func TestHandleRespondsWithFailWhenGetPostsFail(t *testing.T) {
	accessTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token": "abc123"}`))
	}))
	postsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	h, err := NewHandler("clientId", "clientSecret", accessTokenServer.URL, postsServer.URL)
	if err != nil {
		t.Errorf("failed with error: %+v", err)
	}

	req := events.APIGatewayProxyRequest{}
	expected := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
	resp, err := h.Handle(req)
	if err == nil {
		t.Errorf("got no expected error")
	}
	if resp.StatusCode != expected.StatusCode {
		t.Errorf("request did not fail, status code %d", resp.StatusCode)
	}
}

func TestHandleRespondsWithCodeAndLocationWhenOK(t *testing.T) {

	accessTokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token": "abc123"}`))
	}))
	postsFixture, err := os.ReadFile("../../test/posts-response.json")
	if err != nil {
		t.Errorf("failed reading posts fixture")
	}

	postsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(postsFixture)
	}))

	h, err := NewHandler("clientId", "clientSecret", accessTokenServer.URL, postsServer.URL)
	if err != nil {
		t.Errorf("failed with error: %+v", err)
	}

	req := events.APIGatewayProxyRequest{}
	expected := events.APIGatewayProxyResponse{StatusCode: http.StatusTemporaryRedirect}
	resp, err := h.Handle(req)
	if err != nil {
		t.Errorf("got unexpected error %+v", err)
	}
	if resp.StatusCode != expected.StatusCode {
		t.Errorf("request failed, status code %d", resp.StatusCode)
	}

	hasLocationHeader := false
	for key := range resp.Headers {
		if key == "Location" {
			hasLocationHeader = true
		}
	}
	if hasLocationHeader != true {
		t.Errorf("missing location header")
	}
}
