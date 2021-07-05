package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/anodot/anodot-common/pkg/metrics3"
)

const (
	DATA_TOKEN = "your-data-token"
	API_TOKEN  = "your-api-token"
)

var schema metrics3.AnodotMetricsSchema = metrics3.AnodotMetricsSchema{
	Name: "schema_test",
	Dimensions: []string{
		"GEO",
		"OS",
	},
	Measurements: map[string]metrics3.MeasurmentBase{
		"req_num": metrics3.MeasurmentBase{
			Aggregation: "average",
			CountBy:     "none",
		},
		"req_latency": metrics3.MeasurmentBase{
			Aggregation: "average",
			CountBy:     "none",
			Units:       "ms",
		},
	},
	MissingDimPolicy: &metrics3.DimensionPolicy{
		Action: "fail",
	},
}

var metrics metrics3.AnodotMetrics30 = metrics3.AnodotMetrics30{
	Dimensions: map[string]string{
		"GEO": "Kyiv",
		"OS":  "Macos",
	},
	Measurements: map[string]float64{
		"req_num":     12000,
		"req_latency": 50,
	},
	Tags: map[string][]string{
		"some_tag": []string{
			"value1", "value2",
		},
	},
	Timestamp: metrics3.AnodotTimestamp{time.Now()},
}

func main() {
	var schemaId string

	url, _ := url.Parse("https://app.anodot.com")

	dataToken, err := metrics3.NewAnoToken(DATA_TOKEN, metrics3.DataToken)
	if err != nil {
		panic(err)
	}
	apiToken, err := metrics3.NewAnoToken(API_TOKEN, metrics3.ApiToken)
	if err != nil {
		panic(err)
	}

	client, err := metrics3.NewAnodot30Client(*url, apiToken, nil)
	if err != nil {
		panic(err)
	}
	respCreate, err := client.CreateSchema(schema)
	if err != nil {
		panic(err)
	}

	if respCreate.HasErrors() {
		fmt.Println(respCreate.ErrorMessage())

	}

	respGetschemas, err := client.GetSchemas()
	if err != nil {
		panic(err)
	}

	if respGetschemas.HasErrors() {
		fmt.Println(respGetschemas.ErrorMessage())
	}

	for _, s := range respGetschemas.Schemas {
		if s.Name == "schema_test" {
			schemaId = s.Id
		}
	}

	// Change toke to data token for metrics related requests
	client.Token = dataToken

	// Set schema id for metrics
	metrics.SchemaId = schemaId

	respSubmit, err := client.SubmitMetrics([]metrics3.AnodotMetrics30{metrics})
	if err != nil {
		panic(err)
	}

	if respSubmit.HasErrors() {
		fmt.Println(respSubmit.ErrorMessage())
	}

	// Get next hour to close data bucket: https://docs.anodot.com/#send-stream-watermark
	nextHour := time.Now().Add(time.Hour).Round(time.Hour)

	respWatermark, err := client.SubmitWatermark(schemaId, metrics3.AnodotTimestamp{nextHour})
	if err != nil {
		panic(err)
	}

	if respWatermark.HasErrors() {
		fmt.Println(respWatermark.ErrorMessage())
	}

	respDeleteSchema, err := client.DeleteSchema(schemaId)
	if respDeleteSchema.HasErrors() {
		fmt.Println(respDeleteSchema.ErrorMessage())
	}

	if respWatermark.HasErrors() {
		fmt.Println(respWatermark.ErrorMessage())
	}
}
