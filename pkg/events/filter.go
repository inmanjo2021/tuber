package events

import (
	"regexp"
	"tuber/pkg/util"
)

type image struct {
	name string
	branch string
}

func filter(e *util.RegistryEvent) (qualified bool, event *image) {
	imageNameRegex := regexp.MustCompile(`us\.gcr\.io\/(.*):`)
	name := imageNameRegex.FindString(e.Tag)
	branchRegex := regexp.MustCompile(`us\.gcr\.io\/.*:(.*)`)
	branch := branchRegex.FindString(e.Tag)
	if name == "tuber" && branch == "master" {
		event := &image{name: name, branch: branch}
		return true, event
	}
	return
}