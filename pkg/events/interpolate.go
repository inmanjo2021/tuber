package events

import (
	"sync"
	"tuber/pkg/dataTemplate"
)

func interpolate(yamls []dataTemplate.Yaml, data map[string]string) (interpolated []dataTemplate.Yaml, err error) {
	var chaml = make(chan processedYaml, 1)
	var wg = sync.WaitGroup{}
	go func() {
		for result := range chaml {
			err = result.err
			interpolated = append(interpolated, result.yaml)
			wg.Done()
		}
	}()
	for _, yaml := range yamls {
		wg.Add(1)
		go func(yaml dataTemplate.Yaml) {
			interpolatedYaml, interpolationErr := yaml.NewInterpolated(data)
			chaml <- processedYaml{yaml: interpolatedYaml, err: interpolationErr}
		}(yaml)
	}
	wg.Wait()
	close(chaml)
	return
}

type processedYaml struct {
	yaml dataTemplate.Yaml
	err  error
}
