package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Interface interface {
	Submit20Metrics(metrics []Anodot20Metric) (AnodotResponse, error)
	DeleteMetrics(expressions ...DeleteExpression) (AnodotResponse, error)
}

type metricsService struct {
	client.AnodotClient
}

func New(c client.AnodotClient) *metricsService {
	return &metricsService{c}
}

type Anodot20Metric struct {
	Properties map[string]string      `json:"properties"`
	Timestamp  common.AnodotTimestamp `json:"timestamp"`
	Value      float64                `json:"value"`
	Tags       map[string]string      `json:"tags"`
}

func (m *Anodot20Metric) MarshalJSON() ([]byte, error) {
	type Alias Anodot20Metric

	encProps := make(map[string]string, len(m.Properties))
	encTags := make(map[string]string, len(m.Tags))

	for k, v := range m.Properties {
		encProps[escape(strings.TrimSpace(k))] = escape(strings.TrimSpace(v))
	}

	for k, v := range m.Tags {
		encTags[escape(strings.TrimSpace(k))] = escape(strings.TrimSpace(v))
	}

	return json.Marshal(&struct {
		Properties map[string]string `json:"properties"`
		Tags       map[string]string `json:"tags"`
		*Alias
	}{
		Properties: encProps,
		Tags:       encTags,
		Alias:      (*Alias)(m),
	})
}

func escape(s string) string {
	result := strings.ReplaceAll(s, ".", "_")
	result = strings.ReplaceAll(result, "=", "_")

	return strings.ReplaceAll(result, " ", "_")
}

type AnodotResponse interface {
	HasErrors() bool
	ErrorMessage() string
	RawResponse() *http.Response
}

// Anodot server response.
// See more at: https://app.swaggerhub.com/apis/Anodot/metrics_protocol_2.0/1.0.0#/ErrorResponse
type CreateResponse struct {
	Errors []struct {
		Description string
		Error       int64
		Index       string
	} `json:"errors"`
	HttpResponse *http.Response `json:"-"`
}

func (r *CreateResponse) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *CreateResponse) ErrorMessage() string {
	return fmt.Sprintf("%+v\n", r.Errors)
}

func (r *CreateResponse) RawResponse() *http.Response {
	return r.HttpResponse
}

type DeleteResponse struct {
	ID         string `json:"id"`
	Validation struct {
		Passed   bool `json:"passed"`
		Failures []struct {
			ID      int    `json:"id"`
			Message string `json:"message"`
		} `json:"failures"`
	} `json:"validation"`
	HttpResponse *http.Response `json:"-"`
}

type DeleteExpression struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (a *DeleteResponse) HasErrors() bool {
	return !a.Validation.Passed
}

func (a *DeleteResponse) ErrorMessage() string {
	return fmt.Sprintf("%+v\n", a.Validation.Failures)
}

func (a *DeleteResponse) RawResponse() *http.Response {
	return a.HttpResponse
}

type Submitter interface {
	SubmitMetrics(metrics []Anodot20Metric) (AnodotResponse, error)
	AnodotURL() *url.URL
}

func (s *metricsService) Submit20Metrics(metrics []Anodot20Metric) (AnodotResponse, error) {
	//s.AnodotURL().Path = "/api/v1/metrics"

	anodotURL := s.AnodotURL()
	q := anodotURL.Query()
	q.Set("token", s.Token())
	q.Set("protocol", "anodot20")

	//s.AnodotURL().RawQuery = q.Encode()

	b, e := json.Marshal(metrics)
	if e != nil {
		return nil, fmt.Errorf("Failed to parse message:" + e.Error())
	}

	r, _ := http.NewRequest(http.MethodPost, anodotURL.String(), bytes.NewBuffer(b))
	r.Header.Add("Content-Type", "application/json")

	resp, err := s.HTTPClient().Do(r)
	anodotResponse := &CreateResponse{HttpResponse: resp}
	if err != nil {
		return anodotResponse, err
	}

	if resp.StatusCode != 200 {
		return anodotResponse, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		return anodotResponse, fmt.Errorf("empty response body")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, anodotResponse)
	if err != nil {
		return anodotResponse, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	if anodotResponse.HasErrors() {
		return anodotResponse, errors.New(anodotResponse.ErrorMessage())
	} else {
		return anodotResponse, nil
	}
}

func (s *metricsService) DeleteMetrics(expressions ...DeleteExpression) (AnodotResponse, error) {
	//s.AnodotURL().Path = "/api/v1/metrics"

	anodotURL := s.AnodotURL()
	q := anodotURL.Query()
	q.Set("token", s.Token())

	//s.AnodotURL().RawQuery = q.Encode()

	deleteStruct := struct {
		Expression []DeleteExpression `json:"expression"`
	}{}
	deleteStruct.Expression = expressions

	b, e := json.Marshal(deleteStruct)
	if e != nil {
		return nil, fmt.Errorf("failed to parse delete expression:" + e.Error())
	}

	r, _ := http.NewRequest(http.MethodDelete, anodotURL.String(), bytes.NewBuffer(b))
	r.Header.Add("Content-Type", "application/json")

	resp, err := s.HTTPClient().Do(r)
	anodotResponse := &DeleteResponse{HttpResponse: resp}
	if err != nil {
		return anodotResponse, err
	}

	statusCode := resp.StatusCode
	if statusCode < 200 && statusCode >= 300 {
		return anodotResponse, fmt.Errorf("http error: %d", statusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, anodotResponse)
	if err != nil {
		return anodotResponse, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	if resp.Body == nil {
		return anodotResponse, fmt.Errorf("empty response body")
	}

	if anodotResponse.HasErrors() {
		return anodotResponse, errors.New(anodotResponse.ErrorMessage())
	} else {
		return anodotResponse, nil
	}
}
