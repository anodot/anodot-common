package api

import (
	"fmt"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/anodot-common/pkg/metrics"
)

type Api struct {
	*client.AnodotClient
	Events  events.Interface
	Metrics metrics.Interface
}

func NewApiClient(c *client.AnodotClient) (*Api, error) {
	if c == nil {
		return nil, fmt.Errorf("AnodotClient should not be nil")
	}

	return &Api{
		AnodotClient: c,
		Events:       events.NewService(c),
		Metrics:      metrics.NewService(c),
	}, nil
}
