package layers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type authResponse struct {
	Token string `json:"token"`
}

// registry struct
type registry struct {
	baseURL  string
	username string
	password string
	scope    string
}

// NewGoogleRegistry creates registry struct
func newGoogleRegistry(googleToken string) *registry {
	return &registry{
		baseURL:  "https://gcr.io",
		username: "_token",
		password: googleToken,
		scope:    "pull",
	}
}

// getRepository returns repo for image
func (r *registry) getRepository(image string) (repo *repository, err error) {
	token, err := r.getToken(image)

	if err != nil {
		return
	}

	repo = &repository{registry: r, image: image, token: token}
	return
}

// getToken using access token, retrieves request token for registry
func (r *registry) getToken(repository string) (token string, err error) {
	requestURL := fmt.Sprintf("%s/v2/token?scope=repository:%s:%s",
		r.baseURL, repository, r.scope)

	var client = &http.Client{}
	var obj = new(authResponse)

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return
	}

	req.SetBasicAuth(r.username, r.password)
	res, err := client.Do(req)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &obj)

	if err != nil {
		return
	}

	if obj.Token == "" {
		err = fmt.Errorf("no token")
		return
	}

	token = obj.Token
	return
}
