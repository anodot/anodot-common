package metrics3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type AnodotToken struct {
	Token string
	Type  string
}

func NewAnoToken(token string, ttype string) (*AnodotToken, error) {
	if len(strings.TrimSpace(token)) == 0 {
		return nil, fmt.Errorf("token can't be blank")
	}
	if ttype != "api" && ttype != "bearer" {
		return nil, fmt.Errorf("token type can be api or bearer")
	}
	return &AnodotToken{token, ttype}, nil
}

type AnodotResponse interface {
	HasErrors() bool
	ErrorMessage() string
	RawResponse() *http.Response
}

type Anodot30Client struct {
	ServerURL *url.URL
	Token     *AnodotToken
	client    *http.Client
}

func NewAnodot30Client(anodotURL url.URL, token *AnodotToken, httpClient *http.Client) (*Anodot30Client, error) {
	if token == nil {
		return nil, fmt.Errorf("anodot token can't be nil")
	}
	submitter := Anodot30Client{Token: token, ServerURL: &anodotURL, client: httpClient}
	if httpClient == nil {
		client := http.Client{Timeout: 30 * time.Second}

		debugHTTP, _ := strconv.ParseBool(os.Getenv("ANODOT_HTTP_DEBUG_ENABLED"))
		if debugHTTP {
			client.Transport = &debugHTTPTransport{r: http.DefaultTransport}
		}
		submitter.client = &client
	}

	return &submitter, nil
}

func (c *Anodot30Client) SubmitMetrics(metrics []AnodotMetrics30) (AnodotResponse, error) {
	if c.Token.Type != "api" {
		return nil, fmt.Errorf("AnodotToken with type api should be provided for metrics submit ")
	}

	sUrl := *c.ServerURL
	sUrl.Path = "api/v1/metrics"

	q := sUrl.Query()
	q.Set("token", c.Token.Token)
	q.Set("protocol", "anodot30")
	sUrl.RawQuery = q.Encode()

	b, e := json.Marshal(metrics)
	if e != nil {
		return nil, fmt.Errorf("Failed to parse schema:" + e.Error())
	}
	r, _ := http.NewRequest(http.MethodPost, sUrl.String(), bytes.NewBuffer(b))
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	anodotResponse := &SubmitMetricsResponse{HttpResponse: resp}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, anodotResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reponse body: %v \n%s", err, string(bodyBytes))
	}

	return anodotResponse, nil
}

func (c *Anodot30Client) CreateSchema(schema AnodotMetricsSchema) (AnodotResponse, error) {
	if c.Token.Type != "bearer" {
		return nil, fmt.Errorf("AnodotToken with type bearer should be provided for schema creation ")
	}

	var bearer = "Bearer " + c.Token.Token
	sUrl := c.ServerURL
	sUrl.Path = "/api/v2/stream-schemas"

	b, e := json.Marshal(schema)
	if e != nil {
		return nil, fmt.Errorf("Failed to parse schema:" + e.Error())
	}

	r, _ := http.NewRequest(http.MethodPost, sUrl.String(), bytes.NewBuffer(b))

	r.Header.Set("Authorization", bearer)
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	anodotResponse := &CreateSchemaResponse{HttpResponse: resp}

	if resp.Body == nil {
		return anodotResponse, fmt.Errorf("empty response body")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode/100 != 2 {

		err = json.Unmarshal(bodyBytes, &anodotResponse.Error)
		if err != nil {
			return nil, fmt.Errorf("failed to parse reponse body: %v \n%s", err, string(bodyBytes))
		}
	}

	return anodotResponse, nil
}

func (c *Anodot30Client) GetSchemas() ([]AnodotMetricsSchema, error) {

	if c.Token.Type != "bearer" {
		return nil, fmt.Errorf("AnodotToken with type bearer should be provided")
	}
	var bearer = "Bearer " + c.Token.Token
	sUrl := c.ServerURL
	sUrl.Path = "/api/v2/stream-schemas/schemas"

	r, _ := http.NewRequest(http.MethodGet, sUrl.String(), nil)

	r.Header.Set("Authorization", bearer)
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	schemasTmp := make([]StreamSchemaWrapper, 0)
	schemas := make([]AnodotMetricsSchema, 0)

	err = json.Unmarshal(bodyBytes, &schemasTmp)
	if err != nil {
		return nil, err
	}

	for _, s := range schemasTmp {
		schemas = append(schemas, s.Wrapper.Schema)
	}

	return schemas, nil
}

type debugHTTPTransport struct {
	r http.RoundTripper
}

func (d *debugHTTPTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequestOut(h, true)
	fmt.Printf("----------------------------------REQUEST----------------------------------\n%s\n", string(dump))
	resp, err := d.r.RoundTrip(h)
	if err != nil {
		fmt.Println("failed to obtain response: ", err.Error())
		return resp, err
	}

	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Printf("----------------------------------RESPONSE----------------------------------\n%s\n----------------------------------\n\n", string(dump))
	return resp, err
}
