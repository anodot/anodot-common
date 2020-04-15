package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Event struct {
	ID          string                  `json:"id,omitempty"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Category    string                  `json:"category"`
	Source      string                  `json:"source"`
	Properties  []EventProperties       `json:"properties"`
	StartDate   common.AnodotTimestamp  `json:"startDate"`
	EndDate     *common.AnodotTimestamp `json:"endDate,omitempty"`
}
type EventProperties struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

type eventsService struct {
	*client.AnodotClient
}

func NewService(c *client.AnodotClient) *eventsService {
	return &eventsService{AnodotClient: c}
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

func (e *eventsService) Create(event Event) (*Event, error) {
	createStruct := struct {
		Event Event `json:"event"`
	}{Event: event}

	request, err := e.NewRequest(http.MethodPost, "/api/v1/user-events", createStruct)
	if err != nil {
		return nil, err
	}
	resp, err := e.Do(request)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var errResp common.ErrorResponse
		err = json.Unmarshal(bodyBytes, &errResp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Anodot sever response: %w. Status code : %d", err, resp.StatusCode)
		}

		if errResp.HasErrors() {
			return nil, errors.New(errResp.ErrorMessage())
		} else {
			return nil, errors.New(string(bodyBytes))
		}
	}

	var newEvent Event
	err = json.Unmarshal(bodyBytes, &newEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Anodot sever response: %w ", err)
	}

	return &newEvent, nil
}

func (e *eventsService) Get(id string) (*Event, error) {
	panic("implement me")
}

func (e *eventsService) Delete(id string) (Event, error) {
	panic("implement me")
}

func (e *eventsService) CreateCategory(name string, imageURL *url.URL) (*EventCategory, error) {
	createStruct := struct {
		Category struct {
			Name     string `json:"name"`
			ImageURL string `json:"imageUrl,omitempty"`
		} `json:"category"`
	}{}

	createStruct.Category.Name = name
	if imageURL != nil {
		createStruct.Category.ImageURL = imageURL.String()
	}

	request, err := e.NewRequest(http.MethodPost, "api/v1/user-events/categories", createStruct)
	if err != nil {
		return nil, err
	}
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

func (e *eventsService) ListCategories() ([]EventCategory, error) {
	request, err := e.NewRequest(http.MethodGet, "/api/v1/user-events/categories", nil)
	if err != nil {
		return nil, err
	}
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

func (e *eventsService) DeleteCategory(id string) {
	panic("implement me")
}

func (e *eventsService) CreateSource(name string, imageURL *url.URL) (*EventSource, error) {
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
	if err != nil {
		return nil, err
	}
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

func (e *eventsService) ListSources() ([]EventSource, error) {
	request, err := e.NewRequest(http.MethodGet, "api/v1/user-events/sources", nil)
	if err != nil {
		return nil, err
	}
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

func (e *eventsService) DeleteSource() {
	panic("implement me")
}
