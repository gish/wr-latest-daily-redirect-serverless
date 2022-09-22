package model

type RedditPost struct {
	IsSelf bool   `json:"is_self"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}
