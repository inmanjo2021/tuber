package core

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/goccy/go-yaml"
)

// ReleaseTubers combines and interpolates with tuber's conventions, and applies them
func ReleaseTubers(tubers []string, app *TuberApp, digest string) ([]byte, error) {
	return ApplyTemplate(app.Name, strings.Join(tubers, "---\n"), tuberData(app, digest))
}

func tuberData(app *TuberApp, digest string) (data map[string]string) {
	return map[string]string{
		"tuberImage": digest,
	}
}

// for testing without committing up to address-service
var testYaml = heredoc.Doc(`
apiVersion: v1
kind: Pod
metadata:
  name: db-migrate
  namespace: address-service
  annotations:
    sidecar.istio.io/inject: "false"
spec:
  restartPolicy: Never
  containers:
  - name: db-migrate
    image: us.gcr.io/freshly-docker/address-service@sha256:531630e7167f57975bc79617271aaae0d635fde8ea8030cc021e796f9b4c8a4f
    command: ["/bin/sh"]
    args: ["-c", "rails db:migrate > /dev/termination-log"]
    terminationMessagePolicy: FallbackToLogsOnError
    envFrom:
      - secretRef:
          name: address-service-env
`)

// success example
var success = heredoc.Doc(`
apiVersion: v1
kind: Pod
metadata:
  name: db-migrate
  namespace: address-service
  annotations:
    sidecar.istio.io/inject: "false"
spec:
  restartPolicy: Never
  containers:
  - name: db-migrate
    image: us.gcr.io/freshly-docker/address-service@sha256:531630e7167f57975bc79617271aaae0d635fde8ea8030cc021e796f9b4c8a4f
    command: ["printenv"]
    args: ["HOSTNAME", "KUBERNETES_PORT"]
    terminationMessagePolicy: FallbackToLogsOnError
    envFrom:
      - secretRef:
          name: address-service-env
`)

type prerelease struct {
	Metadata Metadata
}

type Metadata struct {
	Name string
}

func RunPrerelease(tubersTemp []string, app *TuberApp, digest string) (out []byte, err error) {
	for _, tuber := range tubersTemp {

		// for testing without committing up to address-service
		tubers := []string{testYaml}
		tuber = testYaml

		prereleaser := prerelease{}
		err = yaml.Unmarshal([]byte(tuber), &prereleaser)

		out, err = ReleaseTubers(tubers, app, digest)
		if err != nil {
			return
		}

		out, err = waitForPhase(prereleaser.Metadata.Name, app)
		if err != nil {
			deleteOut, deleteErr := deletePrereleaser(prereleaser.Metadata.Name, app)
			if deleteErr != nil {
				doubleFailOut := []byte(string(out) + "\n also failed delete:" + string(deleteOut))
				doubleFailErr := fmt.Errorf(err.Error() + "\n also failed delete:" + deleteErr.Error())
				return doubleFailOut, doubleFailErr
			}
			return
		}

		return deletePrereleaser(prereleaser.Metadata.Name, app)
	}
	return []byte{}, fmt.Errorf("unhandled prerelease run exit")
}

func waitForPhase(name string, app *TuberApp) ([]byte, error) {
	for {
		time.Sleep(5 * time.Second)
		status, err := checkPhase(name, app)
		if err != nil {
			return []byte{}, err
		}

		switch stringStatus := strings.Trim(string(status), `"`); stringStatus {
		case "Succeeded":
			return []byte{}, nil
		case "Failed":
			message, err := investigateFailure(name, app)
			if err != nil {
				return message, err
			}
			return message, fmt.Errorf(string(message))
		default:
			continue
		}
	}
	return []byte{}, fmt.Errorf("unhandled prerelease phase watch exit")
}

func checkPhase(name string, app *TuberApp) (out []byte, err error) {
	cmd := exec.Command("kubectl", "get", "pod", name, "-n", app.Name, "-o", `go-template="{{.status.phase}}"`)

	out, err = cmd.CombinedOutput()
	return
}

func investigateFailure(name string, app *TuberApp) (out []byte, err error) {
	cmd := exec.Command("kubectl", "get", "pod", name, "-n", app.Name, "-o", `go-template="{{range .status.containerStatuses}}{{.state.terminated.message}}{{end}}"`)

	out, err = cmd.CombinedOutput()
	return
}

func deletePrereleaser(name string, app *TuberApp) (out []byte, err error) {
	cmd := exec.Command("kubectl", "delete", "pod", name, "-n", app.Name)

	out, err = cmd.CombinedOutput()
	return
}
