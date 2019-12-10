package main

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
	"github.com/joho/godotenv"
)

const megabyte = 1_000_000
const maxSize = megabyte * 1

type authResponse struct {
	Token string `json:"token"`
}

// Layer is a layer
type Layer struct {
	Digest string `json:"digest"`
	Size   int32  `json:"size"`
}

type manifest struct {
	Layers []Layer `json:"layers"`
}

type notTuberLayerError struct {
	message string
}

func (e *notTuberLayerError) Error() string { return e.message }

func getToken() *authResponse {
	requestURL := fmt.Sprintf(
		"%s/v2/token?scope=repository:%s:pull",
		os.Getenv("AUTH_BASE"),
		os.Getenv("IMAGE_NAME"),
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", requestURL, nil)

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

	var obj = new(authResponse)
	err = json.Unmarshal(body, &obj)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(obj)
	return obj
}

func getLayers() []Layer {
	token := getToken().Token

	requestURL := fmt.Sprintf(
		"%s/v2/%s/manifests/%s",
		os.Getenv("REGISTRY_BASE"),
		os.Getenv("IMAGE_NAME"),
		os.Getenv("IMAGE_TAG"),
	)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", requestURL, nil)
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

	var obj = new(manifest)
	err = json.Unmarshal(body, &obj)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(obj)
	return obj.Layers
}

// Yaml is a yaml
type Yaml struct {
	content  string
	filename string
}

func downloadLayer(layerObj *Layer) ([]Yaml, error) {
	token := getToken().Token
	layer := layerObj.Digest

	requestURL := fmt.Sprintf(
		"%s/v2/%s/blobs/%s",
		os.Getenv("REGISTRY_BASE"),
		os.Getenv("IMAGE_NAME"),
		layer,
	)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	gzipped, _ := gzip.NewReader(res.Body)
	archive := tar.NewReader(gzipped)
	var yamls []Yaml

	for {
		header, err := archive.Next()

		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Contents of %s:\n", header.Name)

		if !strings.HasPrefix(header.Name, ".tuber") {
			return nil, &notTuberLayerError{"Contains stuff other than .tuber"}
		}

		if !strings.HasSuffix(header.Name, ".yaml") {
			continue
		}

		bytes, _ := ioutil.ReadAll(archive)

		var yaml Yaml
		yaml.filename = header.Name
		yaml.content = string(bytes)

		yamls = append(yamls, yaml)

		fmt.Println()
	}

	return yamls, nil
}

func findLayer() ([]Yaml, error) {
	layers := getLayers()

	for _, layer := range layers {
		if layer.Size > maxSize {
			log.Println("Layer to large, skipping...")
			continue
		}

		yamls, err := downloadLayer(&layer)

		if err != nil {
			if _, ok := err.(*notTuberLayerError); ok {
				continue
			}

			log.Fatal(err)
		}

		return yamls, nil
	}

	return nil, fmt.Errorf("no tuber layer found")
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	yamls, err := findLayer()

	if err != nil {
		log.Fatal(err)
	}

	for _, yaml := range yamls {
		spew.Dump(yaml)
	}
}
