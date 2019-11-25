package metrics

import (
	"github.com/anodot/anodot-common/anodotParser"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	s, err := NewAnodot20Submitter("http://nginx-metrics.ano-dev.com:8080", "ee474d2cbe9ca451487e6d883359c51e", nil)
	if err != nil {
		log.Fatalf("%v", err)
	}

	metrics := make([]anodotParser.AnodotMetric, 0)
	metrics = append(metrics, anodotParser.AnodotMetric{
		Properties: nil,
		Timestamp:  1574682241,
		Value:      1,
		Tags:       nil,
	})
	metrics = append(metrics, anodotParser.AnodotMetric{
		Properties: nil,
		Timestamp:  1574682241,
		Value:      1,
		Tags:       nil,
	})

	_, err = s.SubmitMetrics(metrics)
	if err != nil {
		t.Fatal("failed: ", err)
	}
}
