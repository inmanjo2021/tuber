package k8s

import (
	"encoding/json"
	"strings"
)

type k8sMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
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
func (c *Config) Save(namespace string) (err error) {
	config := c.config
	config.Data = c.Data

	var jsondata []byte
	jsondata, err = json.Marshal(config)

	if err != nil {
		return
	}

	Apply(jsondata, namespace)

	return
}

// GetConfig returns a Config struct with a Data element containing config map entries
func GetConfig(name string, namespace string, kind string) (config *Config, err error) {
	result, err := Get(strings.ToLower(kind), name, namespace)

	if err != nil {
		return
	}
	var k8sc k8sConfig

	if result == nil {
		k8sc = k8sConfig{
			APIVersion: "v1",
			Kind:       kind,
			Data:       map[string]string{},
			Metadata:   k8sMetadata{Name: name, Namespace: namespace},
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
