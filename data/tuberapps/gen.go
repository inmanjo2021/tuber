// The following directive is necessary to make the package coherent:

// +build ignore

// run using `go generate ./...`
// generate directive is in data/base.go

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main() {
	Yamls()
}

func Yamls() {
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	var yamls []os.FileInfo

	for _, file := range dir {
		if strings.HasSuffix(file.Name(), ".yaml") {
			yamls = append(yamls, file)
		}
	}

	for _, yaml := range yamls {
		file, err := ioutil.ReadFile(yaml.Name())
		if err != nil {
			panic(err)
		}
		separated := strings.Split(yaml.Name(), ".yaml")
		name := separated[0]
		f, err := os.Create(name + ".go")
		if err != nil {
			panic(err)
		}
		exportName := strings.Title(name)
		t := template.Must(template.New("").Parse(`// Package data is generated
package data

import(
	"github.com/MakeNowJust/heredoc"
)

// {{ .exportName }} is generated. Returns the default {{ .name }} for a new tuber app
var {{ .exportName }} = TuberYaml{
	Filename: "{{ .fileName }}",
	Contents: {{ .name }}Contents(),
}

func {{ .name }}Contents() string {
	return heredoc.Doc(` + "`" + "\n{{ .contents }}`)\n}"))
		err = t.Execute(f, map[string]string{
			"name":       name,
			"exportName": exportName,
			"contents":   string(file),
			"fileName":   yaml.Name(),
		})
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}
