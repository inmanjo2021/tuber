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

var service = heredoc.Doc(`
	apiVersion: v1
	kind: Service
	metadata:
		name: {{.appName}}-service
		namespace: {{.appName}}
	spec:
		ports:
		- port: 9090
			targetPort: {{.port}}
			name: grpc-{{.appName}}
		selector:
			app: {{.appName}}
`)

func Init(appName string, port string) (err error) {
	if err = createTuberDirectory(); err != nil {
		return
	}

	if err = createDeploymentYAML(appName); err != nil {
		return
	}

	if err = createServiceYAML(appName, port); err != nil {
		return
	}

	return
}

func createTuberDirectory() (err error) {
	if err = os.Mkdir(tuberConfigPath, os.ModePerm); os.IsExist(err) {
		return nil
	}

	return
}

func createDeploymentYAML(appName string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML("deployment.yaml", deployment, templateData)
}

func createServiceYAML(appName string, port string) (err error) {
	templateData := map[string]string{
		"appName": appName,
		"port":    port,
	}

	return writeYAML("service.yaml", service, templateData)
}

func writeYAML(fileName string, templateString string, templateData map[string]string) (err error) {
	tpl, err := template.New("tpl").Parse(templateString)
	if err != nil {
		return
	}

	var buff bytes.Buffer

	if err = tpl.Execute(&buff, templateData); err != nil {
		return
	}

	if err = ioutil.WriteFile(tuberConfigPath+"/"+fileName, buff.Bytes(), 0644); err != nil {
		return
	}

	return
}
