package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/gobuffalo/packr"
)

const tuberConfigPath = ".tuber"

// InitTuberApp creates a bunch of yamls for you
func InitTuberApp(appName string, routePrefix string) (err error) {
	if err = createTuberDirectory(); err != nil {
		return
	}

	box := packr.NewBox("../../yamls")

	deploymentYaml, err := box.FindString("deployment.yaml")
	if err != nil {
		return
	}

	if err = createDeploymentYAML(appName, deploymentYaml); err != nil {
		return
	}

	serviceYaml, err := box.FindString("service.yaml")
	if err != nil {
		return
	}

	if err = createServiceYAML(appName, serviceYaml); err != nil {
		return
	}

	virtualServiceYaml, err := box.FindString("virtual_service.yaml")
	if err != nil {
		return
	}

	if err = createVirtualServiceYAML(appName, routePrefix, virtualServiceYaml); err != nil {
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

func createDeploymentYAML(appName string, fileContent string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML("deployment.yaml", fileContent, templateData)
}

func createServiceYAML(appName string, fileContent string) (err error) {
	templateData := map[string]string{
		"appName": appName,
	}

	return writeYAML("service.yaml", fileContent, templateData)
}

func createVirtualServiceYAML(appName string, routePrefix string, fileContent string) (err error) {
	templateData := map[string]string{
		"appName":     appName,
		"routePrefix": routePrefix,
	}

	return writeYAML("virtual_service.yaml", fileContent, templateData)
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
