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
						- containerPort: 3000
`)

var service = heredoc.Doc(`
	apiVersion: v1
	kind: Service
	metadata:
		name: {{.appName}}
		namespace: {{.appName}}
	spec:
		ports:
		- port: 3000
			name: grpc
		selector:
			app: {{.appName}}
`)

var virtualService = heredoc.Doc(`
	apiVersion: networking.istio.io/v1alpha3
	kind: VirtualService
	metadata:
		name: {{.appName}}-ingress
		namespace: {{.appName}}
	spec:
		hosts:
			- "*"
		gateways:
		- istio-system/tls-gateway
		http:
		- match:
			- uri:
				prefix: {{.routePrefix}}
			route:
			- destination:
					host: {{.appName}}
`)

func Init(appName string, routePrefix string) (err error) {
	if err = createTuberDirectory(); err != nil {
		return
	}

	if err = createDeploymentYAML(appName); err != nil {
		return
	}

	if err = createServiceYAML(appName); err != nil {
		return
	}

	if err = createVirtualServiceYAML(appName, routePrefix); err != nil {
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

func createServiceYAML(appName string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML("service.yaml", service, templateData)
}

func createVirtualServiceYAML(appName string, routePrefix string) (err error) {
	templateData := map[string]string{
		"appName":     appName,
		"routePrefix": routePrefix,
	}

	return writeYAML("virtual_service.yaml", virtualService, templateData)
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
