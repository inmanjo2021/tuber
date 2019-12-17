package release

import (
	"fmt"
	"tuber/pkg/apply"
	"tuber/pkg/layers"
)

// New create or update app in kubernetes
func New(name string, tag string, token string) (out []byte, err error) {
	yamls, err := layers.GetGoogleLayer(name, tag, token)

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
