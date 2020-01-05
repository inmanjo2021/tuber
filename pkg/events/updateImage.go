package events

import (
	"fmt"
	"tuber/pkg/util"

	"github.com/itchyny/gojq"
)

func updateImage(yamls []util.Yaml, event *util.RegistryEvent) (withUpdatedDeployment []util.Yaml, err error) {
	var deploymentFound bool
	for _, yaml := range yamls {
		convertedYaml := convert(yaml)
		deployment, err := isDeployment(convertedYaml)
		if err != nil {
			return nil, err
		}

		if deployment {
			updatedYaml, err := updateDeployment(yaml, convertedYaml, event.Digest)
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

func isDeployment(convertedYaml interface{}) (deployment bool, err error) {
	output, err := runJq(".kind", convertedYaml)
	if err != nil {
		return
	}
	deployment = output == "Deployment"
	return
}

func updateDeployment(originalYaml util.Yaml, convertedYaml interface{}, digest string) (updatedYaml util.Yaml, err error) {
	query := fmt.Sprintf(`.spec.template.spec.containers.[0].image = %s`, digest)
	output, err := runJq(query, convertedYaml)
	if err != nil {
		return
	}
	updatedYaml = util.Yaml{Content: output.(string), Filename: originalYaml.Filename}
	return
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}

	return i
}

func runJq(queryInput string, data interface{}) (v interface{}, err error) {
	query, err := gojq.Parse(queryInput)
	iter := query.Run(data)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			v = ""
			return "", err
		}
		return v, nil
	}

	return
}
