package metrics

import (
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
	Submit20Metrics(metrics []Anodot20Metric) (common.AnodotResponse, error)
	DeleteMetrics(expressions ...DeleteExpression) (common.AnodotResponse, error)
}

type metricsService struct {
	*client.AnodotClient
}

func NewService(c *client.AnodotClient) *metricsService {
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
	SubmitMetrics(metrics []Anodot20Metric) (common.AnodotResponse, error)
	AnodotURL() *url.URL
}

func (s *metricsService) Submit20Metrics(metrics []Anodot20Metric) (common.AnodotResponse, error) {
	u := fmt.Sprintf("/api/v1/metrics?protocol=%s", "anodot20")
	request, err := s.NewRequest(http.MethodPost, u, metrics)
	if err != nil {
		return nil, err
	}

	resp, err := s.Do(request)
	anodotResponse := &common.ErrorResponse{HttpResponse: resp}
	if err != nil {
		return nil, err
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

func (s *metricsService) DeleteMetrics(expressions ...DeleteExpression) (common.AnodotResponse, error) {
	deleteStruct := struct {
		Expression []DeleteExpression `json:"expression"`
	}{}
	deleteStruct.Expression = expressions

	request, err := s.NewRequest(http.MethodDelete, "/api/v1/metrics", deleteStruct)
	if err != nil {
		return nil, err
	}

	resp, err := s.Do(request)
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
