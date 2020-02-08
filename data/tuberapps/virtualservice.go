// Package data is generated
package data

import(
	"github.com/MakeNowJust/heredoc"
)

// Virtualservice is generated. Returns the default virtualservice for a new tuber app
func Virtualservice() TuberYaml {
	return TuberYaml{
		Filename: "virtualservice.yaml",
		Contents: virtualserviceContents(),
	}
}

func virtualserviceContents() string {
	return heredoc.Doc(`
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: '{{.appName}}-ingress'
spec:
  hosts:
  - {{"{{.clusterDefaultHost}}"}}
  gateways:
  - {{"{{.clusterDefaultGateway}}"}}
  http:
  - match:
    - uri:
        prefix: {{.routePrefix}}
    route:
    - destination:
        host: {{.appName}}
`)
}