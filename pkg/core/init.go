package core

import (
	"io/ioutil"
	"os"
	data "tuber/data/tuberapps"
)

const tuberConfigPath = ".tuber"

// InitTuberApp creates a bunch of yamls for you
func InitTuberApp(appName string, routePrefix string) (err error) {
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
	if err = os.Mkdir(".tuber", os.ModePerm); os.IsExist(err) {
		return nil
	}
	return
}

func createDeploymentYAML(appName string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML(data.Deployment, templateData)
}

func createServiceYAML(appName string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML(data.Service, templateData)
}

func createVirtualServiceYAML(appName string, routePrefix string) (err error) {
	templateData := map[string]string{
		"appName":     appName,
		"routePrefix": routePrefix,
	}

	return writeYAML(data.Virtualservice, templateData)
}

func writeYAML(app data.TuberYaml, templateData map[string]string) (err error) {
	interpolated, err := interpolate(string(app.Contents), templateData)

	if err != nil {
		return
	}

	return ioutil.WriteFile(tuberConfigPath+"/"+app.Filename, interpolated, 0644)
}
