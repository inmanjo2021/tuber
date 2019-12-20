package release

import (
	"fmt"
	"tuber/pkg/containers"
	"tuber/pkg/k8s"
	"tuber/pkg/pulp"
)

// New create or update app in kubernetes
func New(app *pulp.TuberApp, token string) (out []byte, err error) {
	yamls, err := containers.GetTuberLayer(app.GetRepositoryLocation(), token)

	if err != nil {
		return
	}

	fmt.Println("Starting Apply for", app.Name, app.Tag)
	out, err = k8s.Apply(yamls)

	if err != nil {
		return
	}

	return
}
