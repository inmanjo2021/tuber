package core

import (
	"tuber/pkg/k8s"

	"github.com/MakeNowJust/heredoc"
)

// CreateTuberApp adds a new tuber app configuration, including namespace,
// role, rolebinding, and a listing in tuber-apps
func CreateTuberApp(appName string, repo string, tag string) (err error) {
	err = k8s.Create(appName, "namespace", appName)
	if err != nil {
		switch err.(type) {
		case k8s.AlreadyExistsError:
			err = nil
		default:
			return
		}
	}

	appRoleTemplate, appRoleData := appRoles(appName)
	err = ApplyTemplate(appName, appRoleTemplate, appRoleData)

	if err != nil {
		return
	}

	err = k8s.CreateEnv(appName)

	if err != nil {
		return
	}

	err = AddAppConfig(appName, repo, tag)

	if err != nil {
		return
	}

	return
}

func appRoles(namespace string) (template string, data map[string]string) {
	template = heredoc.Doc(`
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

	data = map[string]string{
		"Namespace": namespace,
	}

	return
}
