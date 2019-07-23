package shoutouts_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/jboursiquot/shoutouts"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	cases := []struct {
		scenario string
		shoutout *shoutouts.Shoutout
		cw       *mockCloudWatch
	}{
		{
			scenario: "valid",
			shoutout: shoutouts.New(),
			cw:       &mockCloudWatch{},
		},
	}

	for _, c := range cases {
		t.Run(c.scenario, func(t *testing.T) {
			m := shoutouts.NewMetrics(c.cw)
			assert.NoError(t, m.Capture(context.Background(), c.shoutout))
		})
	}
}

type mockCloudWatch struct {
	cloudwatchiface.CloudWatchAPI
}

func (mcw *mockCloudWatch) PutMetricData(input *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	return nil, nil
}
