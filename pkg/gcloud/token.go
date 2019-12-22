package gcloud

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2/google"
)

// GetAccessToken generates a short-lives token
func GetAccessToken() (accessToken string, err error) {
	jsonData, err := ioutil.ReadFile("./credentials.json")

	if err != nil {
		return
	}

	config, err := google.JWTConfigFromJSON(jsonData,
		"https://www.googleapis.com/auth/devstorage.read_only",
	)

	if err != nil {
		return
	}

	ctx := context.Background()
	token, err := config.TokenSource(ctx).Token()

	return token.AccessToken, nil
}
