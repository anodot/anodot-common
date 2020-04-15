package common

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AnodotTimestamp struct {
	time.Time
}

func (t AnodotTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(t.Unix())), nil
}

func (t AnodotTimestamp) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)

	i, err := strconv.ParseInt(strInput, 10, 64)
	if err != nil {
		panic(err)
	}

	t.Time = time.Unix(i, 0)
	return nil
}

type AnodotResponse interface {
	HasErrors() bool
	ErrorMessage() string
	RawResponse() *http.Response
}

// Anodot server response.
// See more at: https://app.swaggerhub.com/apis/Anodot/metrics_protocol_2.0/1.0.0#/ErrorResponse
type ErrorResponse struct {
	Errors []struct {
		Description string
		Error       int64
		Index       string
	} `json:"errors"`
	HttpResponse *http.Response `json:"-"`
}

func (r *ErrorResponse) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ErrorResponse) ErrorMessage() string {
	return fmt.Sprintf("%+v\n", r.Errors)
}

func (r *ErrorResponse) RawResponse() *http.Response {
	return r.HttpResponse
}
