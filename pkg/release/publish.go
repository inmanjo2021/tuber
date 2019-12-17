package release

import (
	"fmt"
	"tuber/pkg/apply"
	"tuber/pkg/layer"
)

// New create or update app in kubernetes
func New(name string, tag string, token string) (out []byte, err error) {
	yamls, err := layer.GetGoogleLayer(name, tag, token)

	if err != nil {
		return
	}

	fmt.Println("Starting Apply for", name, tag)
	out, err = apply.Apply(yamls)

	if err != nil {
		return
	}

	return
}
