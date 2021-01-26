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

// Registry contains the data needed to successfully obtain information about a repository
type Registry struct {
	baseURL  string
	scope    string
	username string
	password string
	client   *http.Client
}

// NewRegistry returns a pointer to a newly created registry
func NewRegistry(domain string, password string, client *http.Client) *Registry {
	if client == nil {
		client = &http.Client{}
	}

	return &Registry{
		baseURL:  "https://" + domain,
		scope:    "pull",
		username: "_token",
		password: password,
		client:   client,
	}
}

// getRepository returns repo for image
func (r *Registry) getRepository(path string) (repo *repository, err error) {
	token, err := r.token(path)
	if err != nil {
		return
	}

	repo = &repository{registry: r, path: path, token: token}
	return
}

// token using access token, retrieves request token for registry
func (r *Registry) token(path string) (token string, err error) {
	requestURL := fmt.Sprintf("%s/v2/token?scope=repository:%s:%s", r.baseURL, path, r.scope)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(r.username, r.password)
	res, err := r.client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	fmt.Println(body)
	ar := new(authResponse)
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal token response: %s", err)
	} else if ar.Token == "" {
		return "", fmt.Errorf("no token")
	}

	return ar.Token, nil
}
