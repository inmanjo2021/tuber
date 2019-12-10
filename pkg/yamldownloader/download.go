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

const megabyte = 1_000_000
const maxSize = megabyte * 1

type authResponse struct {
	token string
}

type layer struct {
	digest string
	size   int32
}

type manifest struct {
	layers []layer
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

func getLayers() []layer {
	token := getToken().token

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
	return obj.layers
}

func downloadLayer(layerObj *layer) ([]util.Yaml, error) {
	token := getToken().token
	layer := layerObj.digest

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

// FindLayer should be called DownloadYamls or something
func FindLayer() ([]util.Yaml, error) {
	layers := getLayers()

	for _, layer := range layers {
		if layer.size > maxSize {
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
