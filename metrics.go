package shoutouts

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"

	gkm "github.com/go-kit/kit/metrics/cloudwatch"
)

// NewMetrics returns a new Metrics.
func NewMetrics(cw cloudwatchiface.CloudWatchAPI) *Metrics {
	return &Metrics{cw: cw}
}

// Metrics captures shoutout metrics we care about.
type Metrics struct {
	cw cloudwatchiface.CloudWatchAPI
}

// Capture captures all the metrics we care about for a shoutout.
func (m *Metrics) Capture(ctx context.Context, shoutout *Shoutout) error {
	km := gkm.New(os.Getenv("METRIC_NAMESPACE"), m.cw)
	dims := map[string]string{"TeamID": shoutout.TeamID}

	measures := []string{"Shoutout"}
	for _, measure := range measures {
		c := km.NewCounter(measure)
		c.With(mapToDims(dims)...).Add(1)
	}

	if err := km.Send(); err != nil {
		return err
	}

	return nil
}

func mapToDims(m map[string]string) []string {
	flat := []string{}
	for key, value := range m {
		flat = append(flat, key)
		flat = append(flat, value)
	}
	return flat
}
