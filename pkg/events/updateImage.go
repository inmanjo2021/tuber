package events

import (
	"fmt"
	"sync"
	"tuber/pkg/util"

	"github.com/icza/dyno"
	"gopkg.in/yaml.v2"
)

func updateImage(yamls []util.Yaml, event *util.RegistryEvent) (withUpdatedDeployment []util.Yaml, err error) {
	var deploymentFound bool
	var chaml = make(chan processedYaml, 1)
	var wg = sync.WaitGroup{}
	go func() {
		for result := range chaml {
			err = result.err
			if result.deployment {
				deploymentFound = true
			}
			withUpdatedDeployment = append(withUpdatedDeployment, result.yaml)
			wg.Done()
		}
	}()
	for _, yaml := range yamls {
		wg.Add(1)
		go func(yaml util.Yaml) {
			processYaml(yaml, event, chaml)
		}(yaml)
	}
	wg.Wait()
	close(chaml)
	if deploymentFound != true {
		err = fmt.Errorf("no deployment found")
	}
	return
}

type processedYaml struct {
	yaml       util.Yaml
	deployment bool
	err        error
}

func processYaml(yaml util.Yaml, event *util.RegistryEvent, chaml chan<- processedYaml) {
	parsedYaml, err := parseYaml(yaml.Content)
	if err != nil {
		chaml <- processedYaml{yaml: yaml, err: err}
		return
	}

	deployment, err := isDeployment(parsedYaml)
	if err != nil {
		chaml <- processedYaml{yaml: yaml, err: err}
		return
	}

	if deployment {
		updatedYaml, err := updateDeployment(parsedYaml, event.Digest, yaml.Filename)
		chaml <- processedYaml{yaml: updatedYaml, deployment: true, err: err}
	} else {
		chaml <- processedYaml{yaml: yaml, err: err}
	}
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

func updateDeployment(convertedYaml map[string]interface{}, digest string, filename string) (updatedYaml util.Yaml, err error) {
	err = dyno.Set(convertedYaml, digest, "spec", "template", "spec", "containers", 0, "image")
	if err != nil {
		return
	}
	out, err := yaml.Marshal(convertedYaml)
	if err != nil {
		return
	}
	updatedYaml = util.Yaml{Content: string(out), Filename: filename}
	return
}
