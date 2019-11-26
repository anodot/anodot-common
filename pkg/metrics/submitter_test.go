package metrics

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

//vnekhai: TODO add more tests
func TestSubmitter(t *testing.T) {

	i, err := strconv.ParseInt("1540153181", 10, 64)
	if err != nil {
		panic(err)
	}
	ts := time.Unix(i, 0)

	metrics := make([]Anodot20Metric, 0)
	metric := Anodot20Metric{Properties: map[string]string{"what": "test2", "target_type": "gauge", "source": "gotest"}, Timestamp: AnodotTimestamp{ts}, Value: 1, Tags: map[string]string{}}
	metrics = append(metrics, metric)
	metrics = append(metrics, metric)

	mocksubmitter := MockSubmitter{f: func(metrics []Anodot20Metric) {
		b, e := json.Marshal(metrics)
		if e != nil {
			t.Fatal("Failed to convert metrics to json", e)
		}

		excpedtedJonsOutput := `[{"properties":{"source":"gotest","target_type":"gauge","what":"test2"},"timestamp":1540153181,"value":1,"tags":{}},{"properties":{"source":"gotest","target_type":"gauge","what":"test2"},"timestamp":1540153181,"value":1,"tags":{}}]`
		actualJson := string(b)
		if actualJson != excpedtedJonsOutput {
			t.Fatalf("expected metrics json: %v, \n got: %v", excpedtedJonsOutput, actualJson)
		}
	}}
	mocksubmitter.SubmitMetrics(metrics)
}

type MockSubmitter struct {
	f func(metrics []Anodot20Metric)
}

func (s *MockSubmitter) SubmitMetrics(metrics []Anodot20Metric) {
	s.f(metrics)
}
