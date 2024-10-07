package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-chi/chi/v5"
)

// ListContainers returns the list of containers
func (a *API) ListContainers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	containers, err := a.client.ListContainers(docker.ListContainersOptions{
		All:     true,
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, containers)
}

// CreateContainer runs a container
func (a *API) CreateContainer(w http.ResponseWriter, r *http.Request) {
	var c docker.CreateContainerOptions

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		write(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	c.Context = ctx

	create, err := a.client.CreateContainer(c)
	if err != nil {
		if errors.Is(err, docker.ErrContainerAlreadyExists) {
			write(w, http.StatusConflict, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	if err = a.client.StartContainer(create.ID, nil); err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, create)
}

// InspectContainerWithOptions inspects a container
func (a *API) InspectContainerWithOptions(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	c, err := a.client.InspectContainerWithOptions(docker.InspectContainerOptions{
		ID:      id,
		Context: ctx,
	})
	if err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, c)
}

// RemoveContainer removes a container
func (a *API) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := a.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            id,
		RemoveVolumes: false,
		Force:         false,
		Context:       ctx,
	}); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container deleted"})
}

// StartContainer starts a container
func (a *API) StartContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.client.StartContainer(id, nil); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		if err.Error() == (&docker.ContainerAlreadyRunning{ID: id}).Error() {
			write(w, http.StatusConflict, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container started"})
}

// StopContainer stops a container
func (a *API) StopContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.client.StopContainer(id, 0); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		if err.Error() == (&docker.ContainerNotRunning{ID: id}).Error() {
			write(w, http.StatusConflict, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container stopped"})
}

// ContainerLogs returns the logs of a container
func (a *API) ContainerLogs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var buf bytes.Buffer
	err := a.client.Logs(docker.LogsOptions{
		Context:      ctx,
		Container:    id,
		Stdout:       true,
		Stderr:       true,
		OutputStream: &buf,
		ErrorStream:  &buf,
		Timestamps:   true,
	})
	if err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	var logs []string
	for _, line := range bytes.Split(buf.Bytes(), []byte("\n")) {
		if len(line) > 0 {
			logs = append(logs, string(line))
		}
	}

	write(w, http.StatusOK, logs)
}

// RestartContainer restarts a container
func (a *API) RestartContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.client.RestartContainer(id, 0); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container restarted"})
}

// ExportContainer exports a container
func (a *API) ExportContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	var out = &bytes.Buffer{}
	err := a.client.ExportContainer(docker.ExportContainerOptions{
		ID:           id,
		Context:      ctx,
		OutputStream: out,
	})
	if err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+id+".tar")
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := w.Write(out.Bytes()); err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container exported"})
}

// KillContainer kills a container
func (a *API) KillContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.client.KillContainer(docker.KillContainerOptions{ID: id}); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		if err.Error() == (&docker.ContainerNotRunning{ID: id}).Error() {
			write(w, http.StatusConflict, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container killed"})
}

// PauseContainer pauses a container
func (a *API) PauseContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.client.PauseContainer(id); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container paused"})
}

// PruneContainers prunes containers
func (a *API) PruneContainers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pruned, err := a.client.PruneContainers(docker.PruneContainersOptions{
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, pruned)
}

// UnpauseContainer unpauses a container
func (a *API) UnpauseContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := a.client.UnpauseContainer(id); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container unpaused"})
}

// UpdateContainer updates a container
func (a *API) UpdateContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var c docker.UpdateContainerOptions

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		write(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	c.Context = ctx

	if err := a.client.UpdateContainer(id, c); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container updated"})
}

// ResizeContainerTTY resizes the tty of a container
func (a *API) ResizeContainerTTY(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var c struct{ Height, Width int }

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		write(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	if err := a.client.ResizeContainerTTY(id, c.Height, c.Width); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container resized"})
}

// RenameContainer renames a container
func (a *API) RenameContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var c docker.RenameContainerOptions

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		write(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	c.Context = ctx

	if err := a.client.RenameContainer(c); err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Container renamed"})
}

// TopContainer returns the top of a container
func (a *API) TopContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	top, err := a.client.TopContainer(id, "")
	if err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, top)
}

// WaitContainer waits for a container
func (a *API) WaitContainer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	status, err := a.client.WaitContainer(id)
	if err != nil {
		if err.Error() == (&docker.NoSuchContainer{ID: id}).Error() {
			write(w, http.StatusNotFound, Response{Error: err.Error()})
			return
		}

		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, status)
}
