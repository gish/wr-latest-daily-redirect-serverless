package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"wr-latest-daily-redirect-serverless/pkg/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

var commit string

func main() {

	clientId := os.Getenv("REDDIT_CLIENT_ID")
	if clientId == "" {
		panic("REDDIT_CLIENT_ID not set")
	}

	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	if clientSecret == "" {
		panic("REDDIT_CLIENT_SECRET not set")
	}

	appId := os.Getenv("APP_ID")
	if appId == "" {
		panic("APP_ID not set")
	}

	username := os.Getenv("USERNAME")
	if username == "" {
		panic("USERNAME not set")
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if ok != true {
		panic(fmt.Sprintf("failed getting buildinfo"))
	}

	userAgent := buildUserAgent(buildInfo.GoVersion, appId, commit, username)

	h, err := handler.NewHandler(clientId, clientSecret, "https://www.reddit.com/api/v1/access_token", "https://oauth.reddit.com/r/weightroom", userAgent)
	if err != nil {
		panic(fmt.Sprintf("failed creating new handler: %s", err))
	}

	lambda.Start(h.Handle)
}

func buildUserAgent(platform, appId, version, username string) string {
	return fmt.Sprintf("%s:%s:%s (by /u/%s)", platform, appId, version, username)
}
