package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

// API is the handler for the API
type API struct {
	client *docker.Client
	logger *logrus.Logger
}

// NewApi creates a new API
func NewApi(endpoint string, logger *logrus.Logger) (*API, error) {
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.TraceLevel)
		logger.SetReportCaller(true)
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			PadLevelText:    true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return "", fmt.Sprintf(" %s:%d", frame.File, frame.Line)
			},
		})
	}

	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
		logger.Info("Docker endpoint is not set, using default value")
	}

	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &API{client, logger}, err
}

// Router returns the router for the API
func (a *API) Router() func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(a.headersMiddleware)

		r.Route("/containers", func(r chi.Router) {
			r.Get("/", a.ListContainers)        // get the list of containers
			r.Post("/", a.CreateContainer)      // create a container
			r.Post("/prune", a.PruneContainers) // prune containers

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", a.InspectContainerWithOptions) // inspect a container
				r.Get("/logs", a.ContainerLogs)           // get the logs of a container
				r.Get("/stop", a.StopContainer)           // stop a container
				r.Get("/start", a.StartContainer)         // start a container
				r.Get("/restart", a.RestartContainer)     // restart a container
				r.Get("/pause", a.PauseContainer)         // pause a container
				r.Get("/unpause", a.UnpauseContainer)     // unpause a container
				r.Get("/kill", a.KillContainer)           // kill a container
				r.Get("/export", a.ExportContainer)       // export a container
				r.Get("/top", a.TopContainer)             // get the top of a container
				r.Get("/wait", a.WaitContainer)           // wait for a container
				r.Post("/rename", a.RenameContainer)      // rename a container
				r.Post("/update", a.UpdateContainer)      // update a container
				r.Post("/resize", a.ResizeContainerTTY)   // resize a container
				r.Delete("/", a.RemoveContainer)          // remove a container
			})
		})

		r.Route("/networks", func(r chi.Router) {
			r.Get("/", a.GetNetworks) // get the list of networks
		})

		r.Route("/volumes", func(r chi.Router) {
			r.Get("/", a.GetVolumes) // get the list of volumes
		})

		r.Route("/images", func(r chi.Router) {
			r.Get("/", a.ListImages)             // get the list of images
			r.Get("/search", a.SearchImages)     // search images
			r.Get("/searchEx", a.SearchImagesEx) // search images
			r.Get("/export", a.ExportImages)     // export images
			r.Post("/prune", a.PruneImages)      // prune images

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", a.InspectImage)                   // inspect an image
				r.Get("/history", a.ImageHistory)            // get the history of an image
				r.Get("/export", a.ExportImage)              // export an image
				r.Get("/import", a.ImportImage)              // import an image
				r.Get("/build", a.BuildImage)                // build an image
				r.Post("/tag", a.TagImage)                   // tag an image
				r.Post("/push", a.PushImage)                 // push an image
				r.Post("/pull", a.PullImage)                 // pull an image
				r.Post("/load", a.LoadImage)                 // load an image
				r.Delete("/extended", a.RemoveImageExtended) // remove an image with options
				r.Delete("/", a.RemoveImage)                 // remove an image
			})
		})
	}
}

func (a *API) headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Response struct {
	Message string `json:",omitempty"`
	Error   string `json:",omitempty"`
}

// write writes the response
func write(w http.ResponseWriter, statusCode int, data interface{}) {
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	if data == nil {
		if statusCode == http.StatusOK {
			data = Response{Message: http.StatusText(statusCode)}
		} else {
			data = Response{Error: http.StatusText(statusCode)}
		}
	}

	_ = json.NewEncoder(w).Encode(data)
}
