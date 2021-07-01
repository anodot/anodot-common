package metrics3

import (
	"fmt"
	"net/http"
)

type Measurment struct {
	Base MeasurmentBase
	Name string
}

type MeasurmentBase struct {
	Aggregation string `json:"aggregation"`
	CountBy     string `json:"countBy"`
	Units       string `json:"units,omitempty"`
}

type DimensionPolicy struct {
	Action string `json:"action"`
	Fill   string `json:"fill,omitempty"`
}

type AnodotMetricsSchema struct {
	Id               string                    `json:"id,omitempty"`
	Dimensions       []string                  `json:"dimensions"`
	Measurements     map[string]MeasurmentBase `json:"measurements"`
	MissingDimPolicy *DimensionPolicy          `json:"missingDimPolicy,omitempty"`
	Name             string                    `json:"name"`
	Version          string                    `json:"version,omitempty"`
}

type StreamSchemaWrapper struct {
	Wrapper struct {
		Schema AnodotMetricsSchema `json:"schema"`
	} `json:"streamSchemaWrapper"`
}

type CreateSchemaResponse struct {
	SchemaId *string
	Error    *struct {
		Status        int    `json:"status"`
		Name          string `json:"name"`
		AndtErrorCode int    `json:"andtErrorCode"`
		Path          string `json:"path"`
	}
	HttpResponse *http.Response `json:"-"`
}

func (r *CreateSchemaResponse) HasErrors() bool {
	return r.Error != nil
}

func (r *CreateSchemaResponse) ErrorMessage() string {
	return fmt.Sprintf("%+v\n", r.Error)
}

func (r *CreateSchemaResponse) RawResponse() *http.Response {
	return r.HttpResponse
}

type GetSchemaResponse struct {
	Schemas []AnodotMetricsSchema
	Error   *struct {
		Status        int    `json:"status"`
		Name          string `json:"name"`
		AndtErrorCode int    `json:"andtErrorCode"`
		Path          string `json:"path"`
	}
	HttpResponse *http.Response `json:"-"`
}

func (r *GetSchemaResponse) HasErrors() bool {
	return r.Error != nil
}

func (r *GetSchemaResponse) ErrorMessage() string {
	return fmt.Sprintf("%+v\n", r.Error)
}

func (r *GetSchemaResponse) RawResponse() *http.Response {
	return r.HttpResponse
}
