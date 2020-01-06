package events

import (
	"fmt"
	"tuber/pkg/util"

	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

func updateImage(yamls []util.Yaml, event *util.RegistryEvent) (withUpdatedDeployment []util.Yaml, err error) {
	var deploymentFound bool
	for _, yaml := range yamls {
		parsedYaml, err := parseYaml(yaml.Content)
		if err != nil {
			return nil, err
		}

		deployment, err := isDeployment(parsedYaml)
		if err != nil {
			return nil, err
		}

		if deployment {
			updatedYaml, err := updateDeployment(yaml, parsedYaml, event.Digest)
			if err != nil {
				return nil, err
			}
			deploymentFound = true
			withUpdatedDeployment = append(withUpdatedDeployment, updatedYaml)
		} else {
			withUpdatedDeployment = append(withUpdatedDeployment, yaml)
		}
	}
	if deploymentFound != true {
		err = fmt.Errorf("no deployment found")
	}
	return
}

func parseYaml(data string) (parsed map[string]interface{}, err error) {
	err = yaml.Unmarshal([]byte(data), &parsed)
	return
}

func isDeployment(convertedYaml map[string]interface{}) (deployment bool, err error) {
	kind, err := dyno.GetString(convertedYaml, "kind")
	if err != nil {
		return
	}
	deployment = kind == "Deployment"
	return
}

func updateDeployment(originalYaml util.Yaml, convertedYaml map[string]interface{}, digest string) (updatedYaml util.Yaml, err error) {
	err = dyno.Set(convertedYaml, digest, "spec", "template", "spec", "containers", 0, "image")
	if err != nil {
		return
	}
	out, err := yaml.Marshal(convertedYaml)
	if err != nil {
		return
	}
	updatedYaml = util.Yaml{Content: string(out), Filename: originalYaml.Filename}
	return
}
