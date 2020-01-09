package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// CreateFromFile creates a secret based on the contents of a file
func CreateFromFile(path string, namespace string) (dat []byte, err error) {
	dat, err = ioutil.ReadFile(path)
	projectName := "tuber"

	if err != nil {
		return
	}

	str := string(dat)
	filename := filepath.Base(path)
	data := map[string]string{filename: str}
	meta := k8sMetadata{
		Name:      fmt.Sprintf("%s-%s", projectName, filename),
		Namespace: projectName,
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

	Apply(jsondata, namespace)

	return
}
