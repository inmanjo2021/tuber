// Package data is generated
package data

import(
	"github.com/MakeNowJust/heredoc"
)

// Service is generated. Returns the default service for a new tuber app
func Service() TuberYaml {
	return TuberYaml{
		Filename: "service.yaml",
		Contents: serviceContents(),
	}
}

func serviceContents() string {
	return heredoc.Doc(`
apiVersion: v1
kind: Service
metadata:
  name: {{.appName}}
spec:
  ports:
  - port: 3000
    name: grpc
  selector:
    app: {{.appName}}
`)
}