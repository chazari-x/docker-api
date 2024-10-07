package api

import (
	"net/http"
)

// GetNetworks returns the list of networks
func (a *API) GetNetworks(w http.ResponseWriter, r *http.Request) {
	networks, err := a.client.ListNetworks()
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, networks)
}
