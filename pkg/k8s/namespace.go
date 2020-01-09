package k8s

import (
	"bytes"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/MakeNowJust/heredoc"
)

// CreateNamespace create a new namespace in kubernetes
func CreateNamespace(namespace string) (err error) {
	cmd := exec.Command("kubectl", "create", "namespace", namespace)

	out, err := cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() != 0 {
		err = fmt.Errorf(string(out))
	}

	return
}

type templater struct {
	Namespace string
}

func ApplyTemplate(namespace string, templatestring string, params map[string]string) (out []byte, err error) {
	tpl, err := template.New("tpl").Parse(templatestring)

	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, params)

	if err != nil {
		return
	}

	out, err = Apply(buf.Bytes(), namespace)

	return
}

// BindNamespace create a new namespace in kubernetes
func BindNamespace(namespace string) ([]byte, error) {
	templatestring := heredoc.Doc(`
		---
		kind: Role
		apiVersion: rbac.authorization.k8s.io/v1beta1
		metadata:
		  name: tuber-admin
		  namespace: {{ .Namespace }}
		rules:
		- apiGroups:
		  - '*'
		  resources:
		  - '*'
		  verbs:
		  - '*'
		---
		kind: RoleBinding
		apiVersion: rbac.authorization.k8s.io/v1beta1
		metadata:
		  name: tuber-admin
		  namespace: {{ .Namespace }}
		roleRef:
		  apiGroup: rbac.authorization.k8s.io
		  kind: Role
		  name: tuber-admin
		subjects:
		- kind: ServiceAccount
		  name: default
		  namespace: tuber
	`)

	params := map[string]string{
		"Namespace": namespace,
	}

	return ApplyTemplate(namespace, templatestring, params)
}
