package k8s

import (
	"encoding/json"
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

// Config represents the editable portion of a configmap
type Config struct {
	config *k8sConfig
	Data   map[string]string
}

// Save persists updates to a configmap to k8s
func (c *Config) Save() (err error) {
	config := c.config
	config.Data = c.Data

	var jsondata []byte
	jsondata, err = json.Marshal(config)

	if err != nil {
		return
	}

	write(jsondata)

	return

}

// GetConfig returns a Config struct with a Data element containing config map entries
func GetConfig(name string) (config *Config, err error) {
	result, err := Get("configmap", name)

	if err != nil {
		return
	}
	var k8sc k8sConfig

	if result == nil {
		k8sc = k8sConfig{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Data:       map[string]string{},
			Metadata:   k8sMetadata{Name: name},
		}
	} else {
		json.Unmarshal(result, &k8sc)
	}

	config = &Config{
		config: &k8sc,
		Data:   k8sc.Data,
	}
	return
}
