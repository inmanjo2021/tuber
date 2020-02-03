package k8s

// CreateNamespace create a new namespace in kubernetes
func CreateNamespace(namespace string) (err error) {
	_, err = Create(namespace, "namespace")
	return
}
