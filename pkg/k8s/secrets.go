package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"tuber/pkg/apply"
)

type k8sMetadata struct {
	Name string `json:"name"`
}

type k8sConfig struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   k8sMetadata       `json:"metadata"`
	Type       string            `json:"type,omitempty"`
	Data       map[string]string `json:"data"`
	StringData map[string]string `json:"stringData,omitempty"`
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
	meta := k8sMetadata{
		Name: fmt.Sprintf("%s-%s", projectName, filename),
	}

	config := k8sConfig{
		APIVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		StringData: data,
		Metadata:   meta,
	}

	var jsondata []byte
	jsondata, err = json.Marshal(config)

	if err != nil {
		return
	}

	apply.Write(jsondata)

	return
}
