package containers

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

	"tuber/pkg/util"
)

type RepositoryLocation struct {
	Host string
	Path string
	Tag  string
}

type RepositoryContainer interface {
	GetRepositoryLocation() RepositoryLocation
}

const megabyte = 1_000_000
const maxSize = megabyte * 1

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

type notTuberLayerError struct {
	message string
}

func (n *notTuberLayerError) Error() string {
	return n.message
}

type Manifest struct {
	Digest string
	Layers []layer    `json:"layers"`
	Errors []apiError `json:"errors"`
}

type Repository struct {
	registry *Registry
	path     string
	token    string
}

// GetTuberLayer downloads yamls for an image
func GetTuberLayer(location RepositoryLocation, password string) (yamls []util.Yaml, err error) {
	registry := NewRegistry(location.Host, password)
	repository, err := registry.GetRepository(location.Path)
	if err != nil {
		return
	}

	yamls, err = repository.findLayer(location.Tag)
	return
}

func (r *Repository) GetManifest(tag string) (manifest Manifest, err error) {
	requestURL := fmt.Sprintf(
		"%s/v2/%s/manifests/%s",
		r.registry.baseURL,
		r.path,
		tag,
	)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	res, err := client.Do(req)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	digest := res.Header["Docker-Content-Digest"][0]

	manifest = Manifest{Digest: digest}
	err = json.Unmarshal(body, &manifest)

	if err != nil {
		return
	}

	if len(manifest.Errors) > 0 {
		err = manifest.Errors[0]
		return
	}

	return manifest, nil
}

func (r *Repository) getLayers(tag string) (layers []layer, err error) {
	manifest, err := r.GetManifest(tag)

	if err != nil {
		return
	}

	return manifest.Layers, nil
}

func (r *Repository) downloadLayer(layerObj *layer) (yamls []util.Yaml, err error) {
	layer := layerObj.Digest

	requestURL := fmt.Sprintf(
		"%s/v2/%s/blobs/%s",
		r.registry.baseURL,
		r.path,
		layer,
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))

	res, err := client.Do(req)

	if err != nil {
		return
	}

	yamls, err = convertResponse(res)
	return
}

func convertResponse(response *http.Response) (yamls []util.Yaml, err error) {
	gzipped, err := gzip.NewReader(response.Body)

	if err != nil {
		return
	}

	archive := tar.NewReader(gzipped)

	for {
		var header *tar.Header
		header, err = archive.Next()

		if err == io.EOF {
			err = nil
			return
		}

		if err != nil {
			return
		}

		if !strings.HasPrefix(header.Name, ".tuber") {
			err = &notTuberLayerError{"contains stuff other than .tuber"}
			return
		}

		if !strings.HasSuffix(header.Name, ".yaml") {
			continue
		}

		var bytes []byte
		bytes, err = ioutil.ReadAll(archive)

		if err != nil {
			return
		}

		yaml := util.Yaml{Filename: header.Name, Content: string(bytes)}

		yamls = append(yamls, yaml)
	}
	return
}

// findLayer finds the .tuber layer containing deploy info for Tuber
func (r *Repository) findLayer(tag string) (yamls []util.Yaml, err error) {
	layers, err := r.getLayers(tag)

	if err != nil {
		return
	}

	for _, layer := range layers {
		if layer.Size > maxSize {
			log.Println("Layer too large, skipping...")
			continue
		}

		yamls, err = r.downloadLayer(&layer)

		if err != nil {
			switch err.(type) {
			case *notTuberLayerError:
				continue
			default:
				return
			}
		}

		return
	}

	err = fmt.Errorf("no tuber layer found")
	return
}
