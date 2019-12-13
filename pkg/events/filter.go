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
	imageName := imageNameRegex.FindString(e.Tag)
	branchRegex := regexp.MustCompile(`us\.gcr\.io\/.*:(.*)`)
	branchName := branchRegex.FindString(e.Tag)
	if imageName == "tuber" && branchName == "master" {
		event := &image{imageName: imageName, branchName: branchName}
		return true, event
	}
	return
}