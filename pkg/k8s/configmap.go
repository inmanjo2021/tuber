package k8s

import (
	"encoding/json"
	"fmt"

	"tuber/pkg/apply"
)

func getConfig(name string) (config k8sConfig, err error) {
	result, err := apply.Get("configmap", name)

	if err != nil {
		return
	}

	if result == nil {
		config = k8sConfig{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Data:       map[string]string{},
			Metadata:   k8sMetadata{Name: name},
		}
	} else {
		json.Unmarshal(result, &config)
	}

	return
}

// AddAppConfig add a new configuration to tuber's config map
func AddAppConfig(name string, repo string, tag string) (err error) {
	config, err := getConfig("tuber-apps")

	if err != nil {
		return
	}

	config.Data[name] = fmt.Sprintf("%s:%s", repo, tag)

	var jsondata []byte
	jsondata, err = json.Marshal(config)

	if err != nil {
		return
	}

	apply.Write(jsondata)

	return
}
