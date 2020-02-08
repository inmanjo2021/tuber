// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates stations.go. It can be invoked by running
// go generate

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
	dir, err := ioutil.ReadDir("data/tuberapps")
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
		file, err := ioutil.ReadFile("data/tuberapps/" + yaml.Name())
		if err != nil {
			panic(err)
		}
		separated := strings.Split(yaml.Name(), ".yaml")
		name := separated[0]
		f, err := os.Create("data/tuberapps/" + name + ".go")
		if err != nil {
			panic(err)
		}
		exportName := strings.Title(name)
		t := template.Must(template.New("").Parse(`package data

import(
	"github.com/MakeNowJust/heredoc"
)

// {{ .exportName }} is generated. Returns the default {{ .name }} for a new tuber app
func {{ .exportName }}() TuberYaml {
	return TuberYaml{
		Filename: "{{ .fileName }}",
		Contents: {{ .name }}Contents(),
	}
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

	f, err := os.Create("data/tuberapps/base.go")
	if err != nil {
		panic(err)
	}
	base := `package data

// TuberYaml is generated. It's a generic representation of a default yaml for new tuber apps
type TuberYaml struct {
	Filename string
	Contents string
}
`
	_, err = f.Write([]byte(base))
	if err != nil {
		panic(err)
	}
}
