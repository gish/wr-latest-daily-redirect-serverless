package daily

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"wr-latest-daily-redirect-serverless/model"
)

type Daily struct {
	ClientId     string
	ClientSecret string
	UserAgent    string
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type PostListingContent struct {
	Data PostListingContentData `json:"data"`
}

type PostListingContentData struct {
	Children []PostListingContentDataChildren `json:"children"`
}

type PostListingContentDataChildren struct {
	Post model.RedditPost `json:"data"`
}

func NewDaily(clientId string, clientSecret string, userAgent string) *Daily {
	daily := &Daily{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		UserAgent:    userAgent,
	}
	return daily
}

func (d *Daily) Posts() (*[]model.RedditPost, error) {
	accessToken, err := d.accessToken("https://www.reddit.com/api/v1/access_token")
	if err != nil {
		return nil, err
	}

	posts, err := d.subRedditPosts("https://oauth.reddit.com/r/weightroom", *accessToken)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (d *Daily) Latest(posts *[]model.RedditPost) (*model.RedditPost, error) {
	titleRegex := regexp.MustCompile("(?i)daily")

	for _, post := range *posts {
		if titleRegex.Match([]byte(post.Title)) && post.IsSelf {
			return &post, nil
		}
	}
	return nil, errors.New("no daily post found")
}

func (d *Daily) accessToken(baseURL string) (*string, error) {
	client := &http.Client{}
	data := url.Values{
		"grant_type": {"client_credentials"},
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(d.ClientId, d.ClientSecret)
	req.Header.Add("User-Agent", d.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed with response %d", resp.StatusCode))
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	a := &AccessTokenResponse{}

	err = json.Unmarshal(content, a)
	if err != nil {
		return nil, err
	}

	if a.AccessToken == "" {
		return nil, errors.New("empty access token")
	}
	return &a.AccessToken, nil
}

func (d *Daily) subRedditPosts(baseURL string, accessToken string) (*[]model.RedditPost, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", accessToken))
	req.Header.Set("User-Agent", d.UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("failed with response %d", resp.StatusCode))
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	postListing := &PostListingContent{}

	err = json.Unmarshal(content, postListing)
	if err != nil {
		return nil, err
	}

	posts := []model.RedditPost{}
	for _, post := range *&postListing.Data.Children {
		posts = append(posts, post.Post)
	}

	return &posts, nil

}
