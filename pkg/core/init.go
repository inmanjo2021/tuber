package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"
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
	tpl, err := template.New("").Parse(string(app.Contents))

	if err != nil {
		return
	}

	var buff bytes.Buffer

	if err = tpl.Execute(&buff, templateData); err != nil {
		return
	}

	if err = ioutil.WriteFile(tuberConfigPath+"/"+app.Filename, buff.Bytes(), 0644); err != nil {
		return
	}

	return
}
