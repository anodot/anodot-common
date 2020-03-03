package events

import (
	"encoding/json"
	"fmt"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Event struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Source      string `json:"source"`
	Properties  []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"properties"`
	StartDate common.AnodotTimestamp  `json:"startDate"`
	EndDate   *common.AnodotTimestamp `json:"endDate,omitempty"`
}

type EventCategory struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`
	ImageURL string `json:"imageUrl"`
	Name     string `json:"name"`
}

type EventSource struct {
	ID       string `json:"id"`
	ImageURL string `json:"imageUrl"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
}

type EventsService struct {
	*client.AnodotClient
}

type Interface interface {
	Create(event Event) (*Event, error)
	Get(id string) (*Event, error)
	Delete(id string) (Event, error)

	CreateCategory(name string, imageURL *url.URL) (*EventCategory, error)
	ListCategories() ([]EventCategory, error)
	DeleteCategory(id string)

	CreateSource(name string, imageURL *url.URL) (*EventSource, error)
	ListSources() ([]EventSource, error)
	DeleteSource()
}

func (e *EventsService) Create(event Event) (*Event, error) {
	createStruct := struct {
		Event Event `json:"event"`
	}{Event: event}

	request, err := e.NewRequest(http.MethodPost, "/api/v1/user-events", createStruct)
	resp, err := e.Do(request)
	if err != nil {
		return nil, err
	}

	var newEvent Event
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &newEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	return &newEvent, nil
}

func (e *EventsService) Get(id string) (*Event, error) {
	panic("implement me")
}

func (e *EventsService) Delete(id string) (Event, error) {
	panic("implement me")
}

func (e *EventsService) CreateCategory(name string, imageURL *url.URL) (*EventCategory, error) {
	createStruct := struct {
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl"`
	}{Name: name}

	if imageURL != nil {
		createStruct.ImageURL = imageURL.String()
	}

	request, err := e.NewRequest(http.MethodPost, "api/v1/user-events/categories", createStruct)
	resp, err := e.Do(request)
	if err != nil {
		return nil, err
	}

	var newCategory EventCategory
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &newCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	return &newCategory, nil
}

func (e *EventsService) ListCategories() ([]EventCategory, error) {
	request, err := e.NewRequest(http.MethodGet, "/api/v1/user-events/categories", nil)
	resp, err := e.Do(request)
	if err != nil {
		return nil, err
	}

	var newCategory []EventCategory
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &newCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	return newCategory, nil
}

func (e *EventsService) DeleteCategory(id string) {
	panic("implement me")
}

func (e *EventsService) CreateSource(name string, imageURL *url.URL) (*EventSource, error) {
	createStruct := struct {
		Source struct {
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl,omitempty"`
		} `json:"source"`
	}{}

	createStruct.Source.Name = name

	if imageURL != nil {
		createStruct.Source.ImageURL = imageURL.String()
	}

	request, err := e.NewRequest(http.MethodPost, "api/v1/user-events/sources", createStruct)
	resp, err := e.Do(request)
	if err != nil {
		return nil, err
	}

	var newSource EventSource
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &newSource)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	return &newSource, nil
}

func (e *EventsService) ListSources() ([]EventSource, error) {
	request, err := e.NewRequest(http.MethodGet, "api/v1/user-events/sources", nil)
	resp, err := e.Do(request)
	if err != nil {
		return nil, err
	}

	var newSource []EventSource
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &newSource)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	return newSource, nil
}

func (e *EventsService) DeleteSource() {
	panic("implement me")
}
