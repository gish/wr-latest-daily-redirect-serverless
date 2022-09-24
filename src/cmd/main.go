package main

import (
	"fmt"
	"log"
	"os"
	"wr-latest-daily-redirect-serverless/pkg/daily"
)

func main() {
	clientId := os.Getenv("REDDIT_CLIENT_ID")
	if clientId == "" {
		panic("REDDIT_CLIENT_ID not set")
	}

	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	if clientSecret == "" {
		panic("REDDIT_CLIENT_SECRET not set")
	}
	d := daily.NewDaily(clientId, clientSecret, "golang:wr-latest-daily-redirect-serverless:1.0.0 (by /u/murrtu)")

	posts, err := d.Posts()
	if err != nil {
		panic(fmt.Sprintf("failed getting posts: %s", err))
	}

	latest, err := d.Latest(posts)
	if err != nil {
		panic(fmt.Sprintf("failed getting latest post: %s", err))
	}
	log.Printf("latest: %+v", latest)
}
