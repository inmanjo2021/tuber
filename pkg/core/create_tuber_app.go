package core

import (
	data "tuber/data/tuberapps"
	"tuber/pkg/k8s"
)

// CreateTuberApp adds a new tuber app configuration, including namespace,
// role, rolebinding, and a listing in tuber-apps
func CreateTuberApp(appName string, repo string, tag string) (err error) {
	namespaceData := map[string]string{
		"namespace": appName,
	}

	for _, yaml := range []data.TuberYaml{data.Namespace, data.Role, data.Rolebinding} {
		err = ApplyTemplate(appName, string(yaml.Contents), namespaceData)
		if err != nil {
			return
		}
	}

	err = k8s.CreateEnv(appName)

	if err != nil {
		return
	}

	if err != nil {
		return
	}

	return AddAppConfig(appName, repo, tag)
}
