package release

import (
	"fmt"
	"tuber/pkg/apply"
	"tuber/pkg/containers"
	"tuber/pkg/k8s"
)

// New create or update app in kubernetes
func New(app *k8s.TuberApp, token string) (out []byte, err error) {
	yamls, err := containers.GetTuberLayer(app, token)

	if err != nil {
		return
	}

	out, err = apply.Apply(yamls)

	if err != nil {
		return
	}

	return
}
