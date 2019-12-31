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

// Registry struct
type Registry struct {
	baseURL  string
	scope    string
	username string
	password string
}

func NewRegistry(domain string, password string) *Registry {
	return &Registry{
		baseURL:  "https://" + domain,
		scope:    "pull",
		username: "_token",
		password: password,
	}
}

// GetRepository returns repo for image
func (r *Registry) GetRepository(path string) (repo *Repository, err error) {
	token, err := r.getToken(path)

	if err != nil {
		return
	}

	repo = &Repository{registry: r, path: path, token: token}
	return
}

// getToken using access token, retrieves request token for registry
func (r *Registry) getToken(path string) (token string, err error) {
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
