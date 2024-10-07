package api

import (
	"context"
	"net/http"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-chi/chi/v5"
)

// ListImages returns the list of images
func (a *API) ListImages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	images, err := a.client.ListImages(docker.ListImagesOptions{
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, images)
}

// ImageHistory returns the history of an image
func (a *API) ImageHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	history, err := a.client.ImageHistory(id)
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, history)
}

// RemoveImage removes an image
func (a *API) RemoveImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := a.client.RemoveImage(id)
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Image removed"})
}

// RemoveImageExtended removes an image with options
func (a *API) RemoveImageExtended(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	force := r.URL.Query().Get("force") == "true"
	noprune := r.URL.Query().Get("noprune") == "true"
	err := a.client.RemoveImageExtended(id, docker.RemoveImageOptions{
		Force:   force,
		NoPrune: noprune,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Image removed"})
}

// InspectImage returns the details of an image
func (a *API) InspectImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	image, err := a.client.InspectImage(id)
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, image)
}

// PushImage pushes an image
func (a *API) PushImage(w http.ResponseWriter, r *http.Request) {
	//id := chi.URLParam(r, "id")
	//err := a.client.PushImage(docker.PushImageOptions{}, docker.AuthConfiguration{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, Response{Message: "Image pushed"})

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}

// PullImage pulls an image
func (a *API) PullImage(w http.ResponseWriter, r *http.Request) {
	//id := chi.URLParam(r, "id")
	//err := a.client.PullImage(docker.PullImageOptions{}, docker.AuthConfiguration{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, Response{Message: "Image pulled"})

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}

// LoadImage loads an image
func (a *API) LoadImage(w http.ResponseWriter, r *http.Request) {
	//err := a.client.LoadImage(docker.LoadImageOptions{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, Response{Message: "Image loaded"})

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}

// ExportImage exports an image
func (a *API) ExportImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := a.client.ExportImage(docker.ExportImageOptions{
		Name:    id,
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Image exported"})
}

// ExportImages exports multiple images
func (a *API) ExportImages(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query()["id"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := a.client.ExportImages(docker.ExportImagesOptions{
		Names:   ids,
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Images exported"})
}

// ImportImage imports an image
func (a *API) ImportImage(w http.ResponseWriter, r *http.Request) {
	//err := a.client.ImportImage(docker.ImportImageOptions{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, Response{Message: "Image imported"})

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}

// BuildImage builds an image
func (a *API) BuildImage(w http.ResponseWriter, r *http.Request) {
	//err := a.client.BuildImage(docker.BuildImageOptions{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, Response{Message: "Image built"})

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}

// TagImage tags an image
func (a *API) TagImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	repo := r.URL.Query().Get("repo")

	tag := r.URL.Query().Get("tag")

	force := r.URL.Query().Get("force") == "true"

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := a.client.TagImage(id, docker.TagImageOptions{
		Repo:    repo,
		Tag:     tag,
		Force:   force,
		Context: ctx,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, Response{Message: "Image tagged"})
}

// SearchImages searches for images
func (a *API) SearchImages(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")

	results, err := a.client.SearchImages(term)
	if err != nil {
		write(w, http.StatusInternalServerError, Response{Error: err.Error()})
		return
	}

	write(w, http.StatusOK, results)
}

// SearchImagesEx searches for images with options
func (a *API) SearchImagesEx(w http.ResponseWriter, r *http.Request) {
	//term := r.URL.Query().Get("term")
	//
	//results, err := a.client.SearchImagesEx(term, docker.AuthConfiguration{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, results)

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}

// PruneImages prunes images
func (a *API) PruneImages(w http.ResponseWriter, r *http.Request) {
	//result, err := a.client.PruneImages(docker.PruneImagesOptions{})
	//if err != nil {
	//	write(w, http.StatusInternalServerError, Response{Error: err.Error()})
	//	return
	//}
	//
	//write(w, http.StatusOK, Response{Message: "Images pruned"})

	write(w, http.StatusNotImplemented, Response{Error: "Not implemented"})
}
