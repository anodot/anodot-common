package metrics

/*func Submit(t *testing.T) {
	os.Setenv("ANODOT_HTTP_DEBUG_ENABLED", "true")

	anodotClient, err := client.NewAnodotClient("http://54.173.224.180", "ffce170d694c831c82ee6747ace5f588", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	service := metricsService{anodotClient}

	metrics := make([]Anodot20Metric, 0)
	metrics = append(metrics, Anodot20Metric{
		Properties: map[string]string{"what": "vova-test", "client": "android"},
		Timestamp:  common.AnodotTimestamp{Time: time.Now()},
		Value:      rand.Float64(),
		Tags:       nil,
	})
	_, err = service.Submit20Metrics(metrics)
	if err != nil {
		t.Fatalf(err.Error())
	}
}*/
