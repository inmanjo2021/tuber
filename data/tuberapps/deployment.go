package data

import(
	"github.com/MakeNowJust/heredoc"
)

// Deployment is generated. Returns the default deployment for a new tuber app
func Deployment() TuberYaml {
	return TuberYaml{
		Filename: "deployment.yaml",
		Contents: deploymentContents(),
	}
}

func deploymentContents() string {
	return heredoc.Doc(`
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{.appName}}
  name: {{.appName}}
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
}