package prometheus

import (
	"fmt"
	"github.com/prometheus/common/model"
	"math"
	"testing"
)

func TestReceiver(t *testing.T) {
	samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value1",
			},
			Timestamp: model.Time(1574693483),
			Value:     333,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value2",
			},
			Timestamp: model.Time(1574693483),
			Value:     8.14,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "test3",
			},
			Timestamp: model.Time(1574693483),
			Value:     6.15,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "pos_inf_value",
			},
			Timestamp: model.Time(1574693483),
			Value:     model.SampleValue(math.Inf(1)),
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "neg_inf_value",
			},
			Timestamp: model.Time(1574693483),
			Value:     model.SampleValue(math.Inf(-1)),
		},
	}

	err, parser := NewAnodotParser(nil, nil)
	if err != nil {
		t.Error(err)
	}

	anodotMetrics := parser.ParsePrometheusRequest(samples)

	if len(anodotMetrics) != 3 {
		t.Fatalf(fmt.Sprintf("Expected number of metris=3. Found=%d", len(anodotMetrics)))
	}

	for _, m := range anodotMetrics {
		switch m.Properties["what"] {
		case "testmetric":
		case "testmetric":

		}
	}

	for _, m := range anodotMetrics {
		_, ok := m.Properties["what"]
		if !ok {
			t.Fatalf(fmt.Sprintf("no what propertry for metric %+v\n", m))
		}
	}
}

func TestFilters(t *testing.T) {
	samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value1",
			},
			Timestamp: model.Time(123456789123),
			Value:     13,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value2",
			},
			Timestamp: model.Time(123456789123),
			Value:     0.99993,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "test3",
			},
			Timestamp: model.Time(123456789123),
			Value:     86.1234,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "pos_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(1)),
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "neg_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(-1)),
		},
	}

	filterOut := `{"test_label":"test_label_value2"}`
	err, parser := NewAnodotParser(nil, &filterOut)

	if err != nil {
		t.Fail()
	}

	metrics := parser.ParsePrometheusRequest(samples)

	if len(metrics) > 3 {
		t.Fail()
	}

}

func TestFilters2(t *testing.T) {
	samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value1",
			},
			Timestamp: model.Time(123456789123),
			Value:     1.11,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value2",
			},
			Timestamp: model.Time(123456789123),
			Value:     2,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "test3",
			},
			Timestamp: model.Time(123456789123),
			Value:     0,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "pos_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(1)),
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "neg_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(-1)),
		},
	}

	filterIn := `{"test_label":"test_label_value2"}`
	err, parser := NewAnodotParser(&filterIn, nil)

	if err != nil {
		t.Fail()
	}

	metrics := parser.ParsePrometheusRequest(samples)

	if len(metrics) > 1 {
		t.Fail()
	}

}

func TestFilters4(t *testing.T) {
	samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"tst_label":           "test_label_value1",
			},
			Timestamp: model.Time(123456789123),
			Value:     1.11,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "testmetric",
				"test_label":          "test_label_value2",
			},
			Timestamp: model.Time(123456789123),
			Value:     2,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "test3",
			},
			Timestamp: model.Time(123456789123),
			Value:     0,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "pos_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(1)),
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "neg_inf_value",
			},
			Timestamp: model.Time(123456789123),
			Value:     model.SampleValue(math.Inf(-1)),
		},
	}

	filterIn := `{"test_label":"test_label_value2","tst_label":"test_label_value1"}`
	err, parser := NewAnodotParser(&filterIn, nil)

	if err != nil {
		t.Fail()
	}

	metrics := parser.ParsePrometheusRequest(samples)

	fmt.Println(metrics)
	if len(metrics) != 2 {
		t.Fail()
	}

}
