package metrics3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AnodotTimestamp struct {
	time.Time
}

func (t AnodotTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(t.Unix())), nil
}

type AnodotMetrics30 struct {
	SchemaId     string              `json:"schemaId"`
	Timestamp    AnodotTimestamp     `json:"timestamp"`
	Dimensions   map[string]string   `json:"dimensions"`
	Measurements map[string]float64  `json:"measurements"`
	Tags         map[string][]string `json:"tags"`
}

func (m *AnodotMetrics30) MarshalJSON() ([]byte, error) {
	type Alias AnodotMetrics30

	dimesnions := make(map[string]string, len(m.Dimensions))
	measurements := make(map[string]float64, len(m.Measurements))

	tags := make(map[string][]string, len(m.Tags))

	for k, v := range m.Dimensions {
		dimesnions[escape(strings.TrimSpace(k))] = escape(strings.TrimSpace(v))
	}
	for k, v := range m.Measurements {
		measurements[escape(strings.TrimSpace(k))] = v
	}

	for k, v := range m.Tags {
		tgs := make([]string, len(v))
		for i, tag := range v {
			tgs[i] = escape(strings.TrimSpace(tag))
		}
		tags[escape(strings.TrimSpace(k))] = tgs
	}

	return json.Marshal(&struct {
		Dimesnions   map[string]string   `json:"dimensions"`
		Measurements map[string]float64  `json:"measurements"`
		Tags         map[string][]string `json:"tags"`
		*Alias
	}{
		Dimesnions:   dimesnions,
		Measurements: measurements,
		Tags:         tags,
		Alias:        (*Alias)(m),
	})
}

type SubmitMetricsResponse struct {
	Errors []struct {
		Description string
		Error       int64
		Index       string
	} `json:"errors"`
	HttpResponse *http.Response `json:"-"`
}

func (r *SubmitMetricsResponse) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *SubmitMetricsResponse) ErrorMessage() string {
	return fmt.Sprintf("%+v\n", r.Errors)
}

func (r *SubmitMetricsResponse) RawResponse() *http.Response {
	return r.HttpResponse
}

func escape(s string) string {
	result := strings.ReplaceAll(s, ".", "_")
	result = strings.ReplaceAll(result, "=", "_")

	return strings.ReplaceAll(result, " ", "_")
}
