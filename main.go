package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func getToken() {
	requestUrl := fmt.Sprintf(
		"%s/v2/token?scope=repository:%s:pull",
		os.Getenv("AUTH_BASE"),
		os.Getenv("IMAGE_NAME"),
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestUrl, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth("_token", os.Getenv("GCLOUD_TOKEN"))
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var obj = new(AuthResponse)
	err = json.Unmarshal(body, &obj)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(obj)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	getToken()
}
