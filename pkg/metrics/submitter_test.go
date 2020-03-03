package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/anodot/anodot-common/pkg/common"
	"reflect"
	"testing"
	"time"
)

func TestSubmitter(t *testing.T) {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	ts, err := time.Parse(layout, str)
	if err != nil {
		t.Fatal(err)
	}

	metrics := make([]Anodot20Metric, 0)
	metric := Anodot20Metric{Properties: map[string]string{"what": "test2", "target_type": "gauge", "source": "gotest"}, Timestamp: common.AnodotTimestamp{ts}, Value: 1, Tags: map[string]string{}}
	metrics = append(metrics, metric)
	metrics = append(metrics, metric)

	mocksubmitter := MockSubmitter{f: func(metrics []Anodot20Metric) {
		b, e := json.Marshal(metrics)
		if e != nil {
			t.Fatal("Failed to convert metrics to json", e)
		}

		excpedtedJonsOutput := `[{"properties":{"source":"gotest","target_type":"gauge","what":"test2"},"timestamp":1415792726,"value":1,"tags":{}},{"properties":{"source":"gotest","target_type":"gauge","what":"test2"},"timestamp":1415792726,"value":1,"tags":{}}]`
		actualJson := string(b)

		equal, err := equalJson(actualJson, excpedtedJonsOutput)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if !equal {
			t.Fatalf("expected metrics json: %v, \n got: %v", excpedtedJonsOutput, actualJson)
		}
	}}
	mocksubmitter.SubmitMetrics(metrics)
}

func TestSpecialChars(t *testing.T) {
	var testData = []struct {
		in          string
		out         string
		description string
	}{
		{"host.port", "host_port", "'.' in name"},
		{"key=value", "key_value", "'=' in name"},
		{"server name", "server_name", "'space' in name"},
		{"name ", "name", "'space' at the end"},
		{"value>S", "value\u003eS", "'>' in name"},
		{"value<X", "value\u003cX", "'<' in name"},
		{"value-1", "value-1", "'-' in name"},
		{"api/v1", "api/v1", "'/' in name"},
		{"host:8080api/v1", "host:8080api/v1", "':' in name"},
		{"a\\m", "a\\\\m", "'\\' in name"},
	}

	t1 := time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)

	for _, v := range testData {

		t.Run(v.description, func(t *testing.T) {
			metric := Anodot20Metric{Properties: map[string]string{"what": v.in, "target_type": "gauge", v.in: "remote_write"}, Timestamp: common.AnodotTimestamp{t1}, Value: 1, Tags: map[string]string{"key": v.in, v.in: "value"}}

			bytes, err := json.Marshal(&metric)
			if err != nil {
				t.Fatalf(err.Error())
			}

			expectedJson := fmt.Sprintf(`{"properties":{"%[1]s":"remote_write","target_type":"gauge","what":"%[1]s"},"tags":{"key":"%[1]s","%[1]s":"value"},"timestamp":1471219200,"value":1}`, v.out)
			equal, err := equalJson(string(bytes), expectedJson)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if !equal {
				t.Fatal(fmt.Sprintf("wrong metrics20 json\n got: %v\n want: %v", string(bytes), expectedJson))
			}
		})
	}
}

func equalJson(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

type MockSubmitter struct {
	f func(metrics []Anodot20Metric)
}

func (s *MockSubmitter) SubmitMetrics(metrics []Anodot20Metric) {
	s.f(metrics)
}
