package k8s

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"tuber/pkg/apply"
)

func getConfig(name string) (config k8sConfig, err error) {
	result, err := apply.Get("configmap", name)

	if err != nil {
		return
	}

	if result == nil {
		config = k8sConfig{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Data:       map[string]string{},
			Metadata:   k8sMetadata{Name: name},
		}
	} else {
		json.Unmarshal(result, &config)
	}

	return
}

// AddAppConfig add a new configuration to tuber's config map
func AddAppConfig(name string, repo string, tag string) (err error) {
	config, err := getConfig("tuber-apps")

	if err != nil {
		return
	}

	config.Data[name] = fmt.Sprintf("%s:%s", repo, tag)

	var jsondata []byte
	jsondata, err = json.Marshal(config)

	if err != nil {
		return
	}

	apply.Write(jsondata)

	return
}

// TuberApp type for tuber app
type TuberApp struct {
	Tag      string
	ImageTag string
	Repo     string
	Name     string
}

type appsCache struct {
	apps      []TuberApp
	expiry    time.Time
}

var cache *appsCache
var mutex sync.Mutex

func (a appsCache) isExpired() bool {
	return cache.expiry.Before(time.Now())
}

func refreshAppsCache(apps []TuberApp) {
	expiry := time.Now().Add(time.Minute * 5)
	cache = &appsCache { apps: apps, expiry: expiry }
}

func getTuberApps() (apps []TuberApp, err error) {
	config, err := getConfig("tuber-apps")

	if err != nil {
		return
	}

	for _, imageTag := range config.Data {
		split := strings.SplitN(imageTag, ":", 2)
		repoSplit := strings.SplitN(split[0], "/", 2)
		repoName := repoSplit[1]

		apps = append(apps, TuberApp{
			Name:     repoName,
			ImageTag: imageTag,
			Tag:      split[1],
			Repo:     split[0],
		})
	}

	return
}

// TuberApps returns a list of tuber apps
func TuberApps() (apps []TuberApp, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if cache == nil || cache.isExpired() {
		apps, err = getTuberApps()

		if err == nil {
			refreshAppsCache(apps)
		}
	}
	apps = cache.apps
	return
}
