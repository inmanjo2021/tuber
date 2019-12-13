package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"tuber/pkg/apply"
)

type metadata struct {
	Name string `json:"name"`
}

type secretsConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   metadata          `json:"metadata"`
	Type       string            `json:"type"`
	Data       map[string]string `json:"stringData"`
}

// CreateFromFile creates a secret based on the contents of a file
func CreateFromFile(path string, mountpoint string) (dat []byte, err error) {
	dat, err = ioutil.ReadFile(path)
	projectName := "tuber"

	if err != nil {
		return
	}

	str := string(dat)
	filename := filepath.Base(path)
	data := map[string]string{filename: str}
	meta := metadata{
		Name: fmt.Sprintf("%s/%s", projectName, filename),
	}

	config := secretsConfig{
		APIVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		Data:       data,
		Metadata:   meta,
	}

	var jsondata []byte
	json.Unmarshal(jsondata, &config)
	apply.Write(jsondata)

	return
}
