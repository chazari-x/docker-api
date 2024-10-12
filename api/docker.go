package api

import (
	"net/http"
)

func (a *API) GetDocker(w http.ResponseWriter, _ *http.Request) {
	docker, err := a.client.Info()
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, docker)
}
