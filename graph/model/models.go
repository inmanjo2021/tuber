package model

import (
	"encoding/json"

	"github.com/freshly/tuber/pkg/db"
)

func (t TuberApp) DBIndexes() (map[string]string, map[string]bool, map[string]int) {
	return map[string]string{
			"name":          t.Name,
			"imageTag":      t.ImageTag,
			"sourceAppName": t.SourceAppName,
		}, map[string]bool{
			"reviewApp": t.ReviewApp,
		}, map[string]int{}
}

func (t TuberApp) DBRoot() string {
	return "apps"
}

func (t TuberApp) DBKey() string {
	return t.Name
}

func (t TuberApp) DBMarshal() ([]byte, error) {
	return json.Marshal(t)
}

func (t TuberApp) DBUnmarshal(data []byte) (db.Model, error) {
	var app TuberApp
	err := json.Unmarshal(data, &app)
	if err != nil {
		return nil, err
	}
	if app.State == nil {
		app.State = &State{}
	}
	return app, nil
}
