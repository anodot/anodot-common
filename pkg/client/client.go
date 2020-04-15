package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type AnodotClient struct {
	serverURL url.URL
	token     string

	httpClient *http.Client
}

func (c *AnodotClient) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *AnodotClient) AnodotURL() url.URL {
	return c.serverURL
}

func (c *AnodotClient) Token() string {
	return c.token
}

func NewAnodotClient(anodotURL string, apiToken string, httpClient *http.Client) (*AnodotClient, error) {

	if len(strings.TrimSpace(apiToken)) == 0 {
		return nil, fmt.Errorf("anodot api token should not be blank")
	}

	parsedUrl, err := url.Parse(anodotURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse anodotURL: %s", err.Error())
	}

	submitter := AnodotClient{token: apiToken, serverURL: *parsedUrl, httpClient: httpClient}
	if httpClient == nil {
		client := http.Client{Timeout: 30 * time.Second}

		debugHTTP, _ := strconv.ParseBool(os.Getenv("ANODOT_HTTP_DEBUG_ENABLED"))
		if debugHTTP {
			client.Transport = &debugHTTPTransport{r: http.DefaultTransport}
		}
		submitter.httpClient = &client
	}

	return &submitter, nil
}

func (c *AnodotClient) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	anodotURL := c.AnodotURL()

	u, err := anodotURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", c.Token())
	u.RawQuery = q.Encode()

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	//TODO change this
	req.Header.Set("User-Agent", "anodot-go-apiclient")

	return req, nil
}

func (c *AnodotClient) Do(request *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient().Do(request)
	if err != nil {
		return nil, err
	}

	statusCode := resp.StatusCode
	if statusCode < 200 && statusCode >= 300 {
		return nil, fmt.Errorf("http error: %d", statusCode)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	return resp, nil
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
