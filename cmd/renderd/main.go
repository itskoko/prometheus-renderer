package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	renderer "github.com/itskoko/prometheus-renderer"
)

const (
	defaultWidth        = 360
	defaultHeight       = defaultWidth
	defaultRangeSeconds = 3600
)

var (
	prometheusAPI = flag.String("u", "http://localhost:9090", "URL of prometheus server")
	listenAddr    = flag.String("l", ":8080", "Address to listen on")
	httpRoot      = flag.String("r", "", "Root path for HTTP endpoints; Use when behind proxy")

	logger = log.With(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)), "caller", log.DefaultCaller)
)

func main() {
	flag.Parse()
	level.Info(logger).Log("msg", "Starting up", "prometheusAPI", *prometheusAPI)
	rdr, err := renderer.New(*prometheusAPI)
	if err != nil {
		level.Error(logger).Log("msg", "Couldn't create renderer", "err", err)
		os.Exit(1)
	}

	http.HandleFunc(*httpRoot+"/graph", func(w http.ResponseWriter, r *http.Request) {
		status, err := render(rdr, w, r)
		if err != nil {
			errStr := err.Error()
			if status == http.StatusInternalServerError {
				errStr = "Internal server error"
				level.Error(logger).Log("msg", errStr, "err", err.Error())
			}
			http.Error(w, errStr, status)
			return
		}
	})

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		level.Error(logger).Log("msg", "Couldn't listen", "err", err)
		os.Exit(1)
	}
}

func render(rdr *renderer.Render, w http.ResponseWriter, r *http.Request) (int, error) {
	query := r.URL.Query().Get("q")
	if query == "" {
		return http.StatusBadRequest, errors.New("Require q parameter")
	}
	hs, err := parseStringDefault(r.URL.Query().Get("h"), defaultHeight)
	if err != nil {
		return http.StatusBadRequest, err
	}
	ws, err := parseStringDefault(r.URL.Query().Get("w"), defaultWidth)
	if err != nil {
		return http.StatusBadRequest, err
	}
	rn, err := parseStringDefault(r.URL.Query().Get("s"), defaultRangeSeconds)
	if err != nil {
		return http.StatusBadRequest, err
	}

	since := time.Duration(rn) * time.Second
	if err := rdr.Render(w, query, since, ws, hs); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func parseStringDefault(in string, dn int) (int, error) {
	if in == "" {
		return dn, nil
	}
	return strconv.Atoi(in)
}
