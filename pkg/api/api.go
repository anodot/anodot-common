package api

import (
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/events"
	"github.com/anodot/anodot-common/pkg/metrics"
)

type Api struct {
	*client.AnodotClient
	Events  events.Interface
	Metrics metrics.Interface
}

func NewApiClient(c *client.AnodotClient) Api {
	return Api{
		AnodotClient: c,
		Events:       &events.EventsService{c},
	}
}
