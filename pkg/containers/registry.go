package containers

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
	scope    string
	username string
	password string
}

func newRegistry(domain string, password string) *registry {
	return &registry{
		baseURL:  "https://" + domain,
		scope:    "pull",
		username: "_token",
		password: password,
	}
}

// getRepository returns repo for image
func (r *registry) getRepository(path string) (repo *repository, err error) {
	token, err := r.getToken(path)

	if err != nil {
		return
	}

	repo = &repository{registry: r, path: path, token: token}
	return
}

// getToken using access token, retrieves request token for registry
func (r *registry) getToken(path string) (token string, err error) {
	requestURL := fmt.Sprintf("%s/v2/token?scope=repository:%s:%s",
		r.baseURL, path, r.scope)

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
