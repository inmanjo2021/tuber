package containers

import (
	"fmt"
	"net/http"
)

type InvalidRegistryResponse struct {
	StatusCode int
	Headers    http.Header
}

func (ir *InvalidRegistryResponse) Error() string {
	return fmt.Sprintf("container registry: returned unexpected status code %d", ir.StatusCode)
}
