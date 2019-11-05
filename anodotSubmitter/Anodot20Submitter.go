package anodotSubmitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/anodot/anodot-common/anodotParser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const PATH = "/api/v1/metrics/"
const CONTENT_TYPE = "application/json"
const METHOD = "POST"
const PROTOCOL = "anodot20"

type Anodot20Submitter struct {
	Url   *url.URL
	Token string

	client *http.Client
}

type AnodotResponse struct {
	Errors []map[string]string `json:"errors"`
}

var (
	anoServerhttpReponses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "anodot_server_http_responses_total",
		Help: "Total number of HTTP reposes of Anodot server",
	}, []string{"response_code"})
)

func NewAnodot20Submitter(urlStr string, token string, httpClient *http.Client) (*Anodot20Submitter, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if len(strings.TrimSpace(token)) == 0 {
		return nil, fmt.Errorf("anodot api token should not be blank")
	}

	submitter := Anodot20Submitter{Token: token, Url: parsedUrl}
	if httpClient == nil {
		submitter.client = &http.Client{Timeout: 30 * time.Second}
	}

	return &submitter, nil
}

func Collectors() []prometheus.Collector {
	collectors := make([]prometheus.Collector, 0)
	collectors = append(collectors, anoServerhttpReponses)

	return collectors
}

func (s *Anodot20Submitter) submitMetrics(metrics []anodotParser.AnodotMetric, u *url.URL, q url.Values) {
	q.Set("protocol", PROTOCOL)

	u.RawQuery = q.Encode()
	urlStr := fmt.Sprintf("%v", u)

	b, e := json.Marshal(metrics)
	if e != nil {
		log.Printf("Failed to parse message:" + e.Error())
		return
	}
	r, _ := http.NewRequest(METHOD, urlStr, bytes.NewBuffer(b))
	r.Header.Add("Content-Type", CONTENT_TYPE)

	resp, err := s.client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	anoServerhttpReponses.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()
	if resp.StatusCode != 200 {
		log.Println("Http Error:", resp.StatusCode)
		return
	}

	if resp.Body == nil {
		fmt.Println("Empty response body")
		//s.Stats.UpdateMeter(remoteStats.REMOTE_HTTP_ERRORS, 1)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var anodotResponse AnodotResponse
	err = json.Unmarshal(bodyBytes, &anodotResponse)
	if err != nil {
		log.Printf("Failed to parse response: " + err.Error())
		//s.Stats.UpdateMeter(remoteStats.REMOTE_HTTP_ERRORS, 1)
	}

	if anodotResponse.Errors != nil {
		fmt.Println(anodotResponse)
	}

}

func (s *Anodot20Submitter) SubmitMetrics(metrics []anodotParser.AnodotMetric) {
	s.Url.Path = PATH
	q := s.Url.Query()
	q.Set("token", s.Token)
	s.submitMetrics(metrics, s.Url, q)

	//s.Stats.UpdateHist(remoteStats.REMOTE_SAMPLES_PER_REQUEST, int64(len(metrics)))
	//s.Stats.UpdateMeter(remoteStats.SUBMITTED_SAMPLES, int64(len(metrics)))
	//s.Stats.UpdateMeter(remoteStats.REMOTE_REQUESTS, 1)
}
