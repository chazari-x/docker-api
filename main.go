package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/chazari-x/docker-api/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	log := logrus.New()

	log.SetLevel(logrus.TraceLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		PadLevelText:    true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf(" %s:%d", frame.File, frame.Line)
		},
	})

	r := chi.NewRouter()

	a, err := api.NewApi("", log)
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
