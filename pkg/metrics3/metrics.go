package metrics3

import (
	"fmt"
	"net/http"
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
