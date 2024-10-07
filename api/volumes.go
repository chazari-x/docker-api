package api

import (
	"context"
	"net/http"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

// GetVolumes returns the list of volumes
func (a *API) GetVolumes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	volumes, err := a.client.ListVolumes(docker.ListVolumesOptions{
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, volumes)
}
