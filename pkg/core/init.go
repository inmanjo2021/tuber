package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/MakeNowJust/heredoc"
)

const tuberConfigPath = ".tuber"

var deployment = heredoc.Doc(`
	apiVersion: apps/v1
	kind: Deployment
	metadata:
		labels:
			app: {{.appName}}
		name: {{.appName}}
		namespace: {{.appName}}
	spec:
		replicas: 1
		selector:
			matchLabels:
				app: {{.appName}}
		template:
			metadata:
				labels:
					app: {{.appName}}
			spec:
				containers:
				- image: {{"{{.tuberImage}}"}}
					name: {{.appName}}
					envFrom:
						- secretRef:
								name: {{.appName}}-env
					ports:
						- containerPort: 80
`)

func CreateTuberDirectory() (err error) {
	if err = os.Mkdir(tuberConfigPath, os.ModePerm); os.IsExist(err) {
		return nil
	}

	return
}

func CreateDeploymentYAML(appName string) (err error) {
	tpl, err := template.New("tpl").Parse(deployment)
	if err != nil {
		return
	}

	var buff bytes.Buffer
	templateData := map[string]string{
		"appName": appName,
	}

	if err = tpl.Execute(&buff, templateData); err != nil {
		return
	}

	if err = ioutil.WriteFile(tuberConfigPath+"/deployment.yaml", buff.Bytes(), 0644); err != nil {
		return
	}

	return
}
