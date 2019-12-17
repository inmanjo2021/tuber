package layers

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

type manifest struct {
	Layers []layer    `json:"layers"`
	Errors []apiError `json:errors`
}

// repository struct for repo
type repository struct {
	registry *registry
	image    string
	token    string
}

// GetGoogleLayer downloads yamls for an image
func GetGoogleLayer(name string, tag string, googleToken string) (yamls []util.Yaml, err error) {
	registry := newGoogleRegistry(googleToken)
	repository, err := registry.getRepository(name)
	if err != nil {
		return
	}

	yamls, err = repository.findLayer(tag)
	return
}

func (r *repository) getLayers(tag string) (layers []layer, err error) {
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
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	var obj = new(manifest)
	err = json.Unmarshal(body, &obj)

	if err != nil {
		return
	}

	if len(obj.Errors) > 0 {
		err = obj.Errors[0]
		return
	}

	layers = obj.Layers
	return
}

func (r *repository) downloadLayer(layerObj *layer) (yamls []util.Yaml, err error) {
	layer := layerObj.Digest

	requestURL := fmt.Sprintf(
		"%s/v2/%s/blobs/%s",
		r.registry.baseURL,
		r.image,
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
			break // End of archive
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

// findLayer finds the .tuber layer containing deploy info for tuber
func (r *repository) findLayer(tag string) (yamls []util.Yaml, err error) {
	layers, err := r.getLayers(tag)

	if err != nil {
		return
	}

	for _, layer := range layers {
		if layer.Size > maxSize {
			log.Println("Layer to large, skipping...")
			continue
		}

		yamls, err = r.downloadLayer(&layer)

		if err != nil {
			if _, ok := err.(*notTuberLayerError); ok {
				continue
			}

			return
		}

		return
	}

	err = fmt.Errorf("no tuber layer found")
	return
}
