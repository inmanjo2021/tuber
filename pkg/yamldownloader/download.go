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
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"tuber/pkg/util"
)

const MEGABYTE = 1_000_000
const MAX_SIZE = MEGABYTE * 1

type AuthResponse struct {
	Token string `json:"token"`
}

type Layer struct {
	Digest string `json:"digest"`
	Size   int32  `json:"size"`
}

type Manifest struct {
	Layers []Layer `json:"layers"`
}

type NotTuberLayerError struct {
	message string
}

func (e *NotTuberLayerError) Error() string { return e.message }

func getToken() *AuthResponse {
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
	return obj
}

func getLayers() []Layer {
	token := getToken().Token

	requestUrl := fmt.Sprintf(
		"%s/v2/%s/manifests/%s",
		os.Getenv("REGISTRY_BASE"),
		os.Getenv("IMAGE_NAME"),
		os.Getenv("IMAGE_TAG"),
	)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var obj = new(Manifest)
	err = json.Unmarshal(body, &obj)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(obj)
	return obj.Layers
}

func DownloadLayer(layerObj *Layer) ([]util.Yaml, error) {
	token := getToken().Token
	layer := layerObj.Digest

	requestUrl := fmt.Sprintf(
		"%s/v2/%s/blobs/%s",
		os.Getenv("REGISTRY_BASE"),
		os.Getenv("IMAGE_NAME"),
		layer,
	)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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
			return nil, &NotTuberLayerError{"Contains stuff other than .tuber"}
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

func FindLayer() ([]util.Yaml, error) {
	layers := getLayers()

	for _, layer := range layers {
		if layer.Size > MAX_SIZE {
			log.Println("Layer to large, skipping...")
			continue
		}

		yamls, err := DownloadLayer(&layer)

		if err != nil {
			if _, ok := err.(*NotTuberLayerError); ok {
				continue
			}

			log.Fatal(err)
		}

		return yamls, nil
	}

	return nil, fmt.Errorf("No tuber layer found.")
}
