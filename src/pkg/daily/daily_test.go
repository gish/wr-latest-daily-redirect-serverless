package daily

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"wr-latest-daily-redirect-serverless/model"
)

func TestAccessTokenReturnsTokenWhenOK(t *testing.T) {
	expected := "abc123"
	expectedUsername := "username"
	expectedPassword := "secretpassword"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		if username != expectedUsername {
			t.Errorf("username not expected: %+v", username)
		}
		if password != expectedPassword {
			t.Errorf("password not expected: %+v", password)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token": "abc123"}`))
	}))
	defer server.Close()
	d := NewDaily(expectedUsername, expectedPassword, "userAgent")
	actual, err := d.accessToken(server.URL)

	if actual == nil {
		t.Errorf("did not get access token [err] [%+v]", err)
	}

	if *actual != expected {
		t.Errorf("did not get access token [actual, expected] [%+v, %+v]", actual, expected)
	}

}

func TestAccessTokenReturnsErrorWhenRequestNotOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"access_token": "abc123"}`))
	}))
	defer server.Close()
	d := NewDaily("username", "password", "userAgent")
	_, err := d.accessToken(server.URL)

	if err.Error() != errors.New("failed with response 500").Error() {
		t.Errorf("did not get error: %+v", err)
	}
}

func TestAccessTokenReturnsErrorWhenParsingResponseNotOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"accessoken": "abc123"}`))
	}))
	defer server.Close()
	d := NewDaily("username", "password", "userAgent")
	_, err := d.accessToken(server.URL)

	if err != nil && err.Error() != errors.New("empty access token").Error() {
		t.Errorf("did not get error: %+v", err)
	}
}

func TestPostsReturnsPostsWhenOK(t *testing.T) {
	postA := model.RedditPost{
		Title:  "Some other thread",
		Url:    "",
		IsSelf: false,
	}
	expected := []model.RedditPost{postA}
	d := NewDaily("id", "secret", "userAgent")
	actual, _ := d.Posts()
	if &expected != actual {
		t.Errorf("posts not returned, [expected, actual] [%+v, %+v]", expected, actual)
	}

}

func SubRedditPostsReturnsErrorWhenRequestFails(t *testing.T) {}

func SubRedditPostsReturnsErrorWhenAccessTokenInvalid(t *testing.T) {}

func SubRedditPostsReturnsPostsWhenOK(t *testing.T) {}

func TestLatestReturnsErrorWhenNoFound(t *testing.T) {
	posts := []model.RedditPost{}
	d := NewDaily("id", "secret", "userAgent")

	latest, err := d.Latest(posts)
	if latest != nil {
		t.Error("latest not nil")
	}
	if err == nil {
		t.Error("error is nil")
	}

}

func TestLatestReturnsFirstFoundDaily(t *testing.T) {
	postA := model.RedditPost{
		Title:  "Some other thread",
		Url:    "",
		IsSelf: false,
	}
	postB := model.RedditPost{
		Title:  "Daily thread",
		Url:    "",
		IsSelf: true,
	}
	postC := model.RedditPost{
		Title:  "External daily thing",
		Url:    "",
		IsSelf: false,
	}
	postD := model.RedditPost{
		Title:  "Daily thread 2",
		Url:    "",
		IsSelf: true,
	}
	posts := []model.RedditPost{postA, postB, postC, postD}
	d := NewDaily("id", "secret", "userAgent")
	latest, _ := d.Latest(posts)

	if *latest != postB {
		t.Errorf("returned post not first found daily, [expected, actual] [%+v, %+v]", postB, latest)
	}
}
