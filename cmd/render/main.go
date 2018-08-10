package main

import (
	"flag"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	renderer "github.com/itskoko/prometheus-renderer"
)

var (
	prometheusAPI = flag.String("u", "http://localhost:9090", "URL of prometheus server")
	filename      = flag.String("f", "out.png", "Path to output file")
	since         = flag.Duration("s", 1*time.Hour, "Graph range")
	width         = flag.Int("w", 800, "Width")
	height        = flag.Int("h", 600, "Height")
	legend        = flag.Bool("l", true, "Show legend")
	logger        = log.With(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)), "caller", log.DefaultCaller)
)

func main() {
	flag.Parse()
	level.Info(logger).Log("msg", "Starting up", "prometheusAPI", *prometheusAPI)
	query := flag.Arg(0)
	if query == "" {
		level.Error(logger).Log("msg", "Query required")
		os.Exit(1)
	}
	f, err := os.Create(*filename)
	if err != nil {
		level.Error(logger).Log("msg", "Couldn't create file", "err", err)
		os.Exit(1)
	}

	r, err := renderer.New(*prometheusAPI)
	if err != nil {
		level.Error(logger).Log("msg", "Couldn't create renderer", "err", err)
		os.Exit(1)
	}
	if err := r.Render(f, query, time.Now().Add(-*since), time.Now(), *width, *height, *legend); err != nil {
		level.Error(logger).Log("msg", "Couldn't render expression", "err", err)
	}
}
