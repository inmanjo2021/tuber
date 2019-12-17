package yamldownloader

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"tuber/pkg/util"
)

const megabyte = 1_000_000
const maxSize = megabyte * 1

type authResponse struct {
	Token string `json:"token"`
}

type layer struct {
	Digest string `json:"digest"`
	Size   int32  `json:"size"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (a apiError) Error() string {
	return a.Message
}

type manifest struct {
	Layers []layer    `json:"layers"`
	Errors []apiError `json:errors`
}

type notTuberLayerError struct {
	message string
}

func (n *notTuberLayerError) Error() string {
	return n.message
}

// Registry struct
type Registry struct {
	baseURL  string
	username string
	password string
}

// NewGoogleRegistry creates Registry struct
func NewGoogleRegistry(googleToken string) *Registry {
	return &Registry{
		baseURL:  "https://us.gcr.io",
		username: "_token",
		password: googleToken,
	}
}

// GetToken using access token, retrieves request token for registry
func (r *Registry) GetToken(repository string, scope string) (token string, err error) {
	requestURL := fmt.Sprintf("%s/v2/token?scope=repository:%s:%s",
		r.baseURL, repository, scope)

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
		return "", fmt.Errorf("no token")
	}

	return obj.Token, nil
}

// Repository struct for repo
type Repository struct {
	registry *Registry
	image    string
	token    string
}

// GetRepository returns repo for image
func (r *Registry) GetRepository(image string, scope string) (repo *Repository, err error) {
	token, err := r.GetToken(image, scope)

	if err != nil {
		return
	}

	return &Repository{
		r,
		image,
		token,
	}, nil

}

func (r *Repository) getLayers(tag string) ([]layer, error) {
	requestURL := fmt.Sprintf(
		"%s/v2/%s/manifests/%s",
		r.registry.baseURL,
		r.image,
		tag,
	)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err

	}

	var obj = new(manifest)
	err = json.Unmarshal(body, &obj)

	if err != nil {
		return nil, err
	}

	if len(obj.Errors) > 0 {
		return nil, obj.Errors[0]
	}

	spew.Dump(obj)
	return obj.Layers, nil
}

func (r *Repository) downloadLayer(layerObj *layer) ([]util.Yaml, error) {
	layer := layerObj.Digest

	requestURL := fmt.Sprintf(
		"%s/v2/%s/blobs/%s",
		r.registry.baseURL,
		r.image,
		layer,
	)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	gzipped, _ := gzip.NewReader(res.Body)
	archive := tar.NewReader(gzipped)
	var yamls []util.Yaml

	for {
		header, err := archive.Next()

		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			log.Fatal(err)
		}

		if !strings.HasPrefix(header.Name, ".tuber") {
			return nil, &notTuberLayerError{"contains stuff other than .tuber"}
		}

		if !strings.HasSuffix(header.Name, ".yaml") {
			continue
		}

		bytes, _ := ioutil.ReadAll(archive)

		var yaml util.Yaml
		yaml.Filename = header.Name
		yaml.Content = string(bytes)

		yamls = append(yamls, yaml)
	}

	return yamls, nil
}

// FindLayer finds the .tuber layer containing deploy info for tuber
func (r *Repository) FindLayer(tag string) ([]util.Yaml, error) {
	layers, err := r.getLayers(tag)

	if err != nil {
		return nil, err
	}

	for _, layer := range layers {
		if layer.Size > maxSize {
			log.Println("Layer to large, skipping...")
			continue
		}

		yamls, err := r.downloadLayer(&layer)

		if err != nil {
			if _, ok := err.(*notTuberLayerError); ok {
				continue
			}

			return nil, err
		}

		return yamls, nil
	}

	return nil, fmt.Errorf("no tuber layer found")
}
