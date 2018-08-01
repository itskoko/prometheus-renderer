package renderer

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type Render struct {
	v1.API
}

func New(prometheusURL string) (*Render, error) {
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		return nil, err
	}
	return &Render{
		API: v1.NewAPI(client),
	}, nil
}

func (r *Render) Render(w io.Writer, query string, since time.Duration, ws, hs int) error {
	plot, err := plot.New()
	if err != nil {
		return err
	}
	plot.Legend.Top = true
	timeseries, err := r.queryRange(query, since)
	if err != nil {
		return err
	}

	if err := plotutil.AddLinePoints(plot, timeseries...); err != nil {
		return err
	}

	plot.Y.Max = plot.Y.Max * 1.2
	c := vgimg.New(vg.Length(ws), vg.Length(hs))
	cpng := vgimg.PngCanvas{Canvas: c}
	cv := draw.New(cpng)
	plot.Draw(cv)

	_, err = cpng.WriteTo(w)
	return err
}

// returns name, plotter.XYer, name1, plotter.XYer ...
func (r *Render) queryRange(query string, since time.Duration) ([]interface{}, error) {
	now := time.Now()

	resp, err := r.API.QueryRange(context.Background(), query, v1.Range{
		Start: now.Add(-since),
		End:   now,
		Step:  60 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	matrix, ok := resp.(model.Matrix)
	if !ok {
		return nil, errors.New("Expected matrix")
	}
	ret := make([]interface{}, 0)
	for _, ss := range matrix {
		pts := make(plotter.XYs, len(ss.Values))
		for i, sample := range ss.Values {
			pts[i].Y = float64(sample.Value)
			pts[i].X = float64(sample.Timestamp.Unix()-now.Unix()) / 60
		}
		ret = append(ret, formatMetric(ss.Metric), pts)
	}
	return ret, nil
}

func formatMetric(m model.Metric) string {
	ls := model.LabelSet(m)
	values := make([]string, len(ls))
	i := 0
	for _, v := range ls {
		values[i] = strings.Replace(string(v), "\n", " ", -1)
		i++
	}
	return strings.Join(values, "|")
}
