package daily

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
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
	d := NewDaily(expectedUsername, expectedPassword, "userAgent", server.URL, "subredditURL")
	actual, err := d.accessToken()

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
	d := NewDaily("username", "password", "userAgent", server.URL, "subredditURL")
	_, err := d.accessToken()

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
	d := NewDaily("username", "password", "userAgent", server.URL, "subredditURL")
	_, err := d.accessToken()

	if err != nil && err.Error() != errors.New("empty access token").Error() {
		t.Errorf("did not get error: %+v", err)
	}
}

func TestSubRedditPostsReturnsErrorWhenRequestFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	d := NewDaily("username", "password", "userAgent", "accessTokenURL", "subredditURL")
	_, err := d.subRedditPosts(server.URL, "accessToken")

	if err.Error() != errors.New("failed with response 500").Error() {
		t.Errorf("did not get error: %+v", err)
	}
}

func TestSubRedditPostsReturnsPostsWhenOK(t *testing.T) {
	postsFixture, err := os.ReadFile("../../test/posts-response.json")
	if err != nil {
		t.Errorf("failed reading posts fixture")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(postsFixture)
	}))
	defer server.Close()
	d := NewDaily("username", "password", "userAgent", "accessTokenURL", "subredditURL")

	posts, err := d.subRedditPosts(server.URL, "accessToken")
	if err != nil {
		t.Errorf("failed getting posts: %+v", err)
	}

	if posts == nil {
		t.Error("posts is nil")
	}

	if len(*posts) == 0 {
		t.Error("posts has no length")
	}

}

func TestLatestReturnsErrorWhenNoFound(t *testing.T) {
	posts := []model.RedditPost{}
	d := NewDaily("username", "password", "userAgent", "accessTokenURL", "subredditURL")

	latest, err := d.Latest(&posts)
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
	d := NewDaily("username", "password", "userAgent", "accessTokenURL", "subredditURL")
	latest, _ := d.Latest(&posts)

	if *latest != postB {
		t.Errorf("returned post not first found daily, [expected, actual] [%+v, %+v]", postB, latest)
	}
}
