package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/markbates/pkger"
)

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

	return writeYAML("deployment.yaml", templateData)
}

func createServiceYAML(appName string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML("service.yaml", templateData)
}

func createVirtualServiceYAML(appName string, routePrefix string) (err error) {
	templateData := map[string]string{
		"appName":     appName,
		"routePrefix": routePrefix,
	}

	return writeYAML("virtual_service.yaml", templateData)
}

func writeYAML(fileName string, templateData map[string]string) (err error) {
	templateDir := pkger.Dir("/yamls")

	templateFile, err := templateDir.Open(fileName)
	if err != nil {
		return
	}
	defer templateFile.Close()

	templateFileBytes, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return
	}

	tpl, err := template.New("tpl").Parse(string(templateFileBytes))
	if err != nil {
		return
	}

	var buff bytes.Buffer

	if err = tpl.Execute(&buff, templateData); err != nil {
		return
	}

	path := "./tuber/" + fileName
	if err = ioutil.WriteFile(path, buff.Bytes(), 0644); err != nil {
		return
	}

	return
}
