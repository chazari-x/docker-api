package main

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"

	"docker/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	log.SetLevel(log.TraceLevel)

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		PadLevelText:    true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf(" %s:%d", frame.File, frame.Line)
		},
	})

	r := chi.NewRouter()

	a, err := api.NewApi("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	r.Route("/api/docker", a.Router())

	log.Trace("Starting server on http://localhost:8080/api/docker")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Error(err)
	}
}
