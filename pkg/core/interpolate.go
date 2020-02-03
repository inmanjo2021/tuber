package core

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
	"tuber/pkg/k8s"

	"github.com/goccy/go-yaml"
)

// ApplyTemplate interpolates and applies a yaml to a given namespace
func ApplyTemplate(namespace string, templatestring string, params map[string]string) (out []byte, err error) {
	tpl, err := template.New("").Parse(templatestring)

	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, params)

	if err != nil {
		return
	}

	out, err = k8s.Apply(buf.Bytes(), namespace)

	return
}

type prerelease struct {
	Kind     string
	Metadata Metadata
}

// Metadata exported for the yaml unmarshaller
type Metadata struct {
	Name string
}

// RunPrerelease takes an array of pods, that are designed to be single use command runners
// that have access to the new code being released.
func RunPrerelease(tubers []string, app *TuberApp, digest string, clusterData *ClusterData) (out []byte, err error) {
	for _, tuber := range tubers {
		prereleaser := prerelease{}
		err = yaml.Unmarshal([]byte(tuber), &prereleaser)
		if err != nil {
			return
		}

		if prereleaser.Kind != "Pod" {
			err = fmt.Errorf("prerelease resources must be Pods, received %s", prereleaser.Kind)
			return
		}

		out, err = ReleaseTubers([]string{tuber}, app, digest, clusterData)
		if err != nil {
			return
		}

		out, err = waitForPhase(prereleaser.Metadata.Name, "pod", app)
		if err != nil {
			deleteOut, deleteErr := k8s.Delete("pod", prereleaser.Metadata.Name, app.Name)
			if deleteErr != nil {
				return deleteOut, fmt.Errorf(err.Error() + "\n also failed delete:" + deleteErr.Error())
			}
			return
		}

		return k8s.Delete("pod", prereleaser.Metadata.Name, app.Name)
	}
	err = fmt.Errorf("unhandled prerelease run exit")
	return
}

func waitForPhase(name string, kind string, app *TuberApp) ([]byte, error) {
	phaseTemplate := fmt.Sprintf(`go-template="%s"`, "{{.status.phase}}")
	failureTemplate := fmt.Sprintf(
		`go-template="%s"`,
		"{{range .status.containerStatuses}}{{.state.terminated.message}}{{end}}",
	)
	timeout := time.Now().Add(time.Minute * 10)

	for {
		if time.Now().After(timeout) {
			return []byte{}, fmt.Errorf("timeout")
		}
		time.Sleep(5 * time.Second)

		status, err := k8s.Get(kind, name, app.Name, "-o", phaseTemplate)
		if err != nil {
			return []byte{}, err
		}

		switch stringStatus := strings.Trim(string(status), `"`); stringStatus {
		case "Succeeded":
			return []byte{}, nil
		case "Failed":
			message, failedRetrieval := k8s.Get(kind, name, app.Name, "-o", failureTemplate)
			if err != nil {
				return message, failedRetrieval
			}
			return message, nil
		default:
			continue
		}
	}
}
