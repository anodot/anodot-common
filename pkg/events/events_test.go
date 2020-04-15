package events

import (
	"encoding/json"
	"fmt"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/common"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestCreateEventSuccessful(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	createBody := `{
  "event": {
    "title": "deployment started on myServer",
    "description": "my description",
    "category": "deployments",
    "source": "chef",
    "properties": [
      {
        "key": "service",
        "value": "myService"
      }
    ],
    "startDate": 1415792726
  }
}`

	mux.HandleFunc("/api/v1/user-events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testTokenPresent(t, r, "test_1234")
		testBody(t, r, createBody)

		fmt.Fprint(w,
			`
{
   "id":"b1229900-1a49-4b9e-a2e4-b9d21240793f",
   "title":"deployment started on myServer",
   "description":"my description",
   "source":"chef",
   "category":"deployments",
   "startDate":1415792726,
   "endDate":null,
   "properties":[
      {
         "key":"service",
         "value":"myService"
      }
   ]
}
`,
		)
	})

	eventsService := eventsService{c}

	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	ts, err := time.Parse(layout, str)
	if err != nil {
		t.Fatal(err)
	}

	eventProperties := []EventProperties{{Key: "service", Value: "myService"}}
	_, err = eventsService.Create(Event{
		Title:       "deployment started on myServer",
		Description: "my description",
		Category:    "deployments",
		Source:      "chef",
		Properties:  eventProperties,
		StartDate:   common.AnodotTimestamp{Time: ts},
		EndDate:     nil,
	})
	if err != nil {
		t.Errorf("eventsService.Create returned error: %v", err)
	}
}

func TestCreateEventError(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	createBody := `{
  "event": {
    "title": "deployment started on myServer",
    "description": "my description",
    "category": "deployments",
    "source": "chef",
    "properties": [
      {
        "key": "service",
        "value": "myService"
      }
    ],
    "startDate": 1415792726
  }
}`

	mux.HandleFunc("/api/v1/user-events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testTokenPresent(t, r, "test_1234")
		testBody(t, r, createBody)

		w.WriteHeader(400)
		fmt.Fprint(w,
			`
{"errors":[{"index":null,"error":3030,"description":"cannot create duplicate event"}]}
`,
		)
	})

	eventsService := eventsService{c}

	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	ts, err := time.Parse(layout, str)
	if err != nil {
		t.Fatal(err)
	}

	eventProperties := []EventProperties{{Key: "service", Value: "myService"}}
	_, err = eventsService.Create(Event{
		Title:       "deployment started on myServer",
		Description: "my description",
		Category:    "deployments",
		Source:      "chef",
		Properties:  eventProperties,
		StartDate:   common.AnodotTimestamp{Time: ts},
		EndDate:     nil,
	})
	if err == nil {
		t.Errorf("eventsService.Create should return error error")
	}

	excpectedErrorMessage := "[{Description:cannot create duplicate event Error:3030 Index:}]\n"
	if err.Error() != excpectedErrorMessage {
		t.Errorf("wrong error message \n got: %v\n want: %v", err.Error(), excpectedErrorMessage)
	}
}

func TestListCategories(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/user-events/categories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testTokenPresent(t, r, "test_1234")
		fmt.Fprint(w,
			`
[
   {
      "id":"0",
      "name":"deployments",
      "imageUrl":"https://s3.amazonaws.com/anodot-images-common/logo-anodot.png",
      "owner":"anodot"
   },
   {
      "id":"1",
      "name":"alerts",
      "imageUrl":"https://s3.amazonaws.com/anodot-images-common/logo-anodot.png",
      "owner":"anodot"
   }
]
`,
		)
	})

	eventsService := eventsService{c}

	artifacts, err := eventsService.ListCategories()
	if err != nil {
		t.Errorf("eventsService.ListCategories returned error: %v", err)
	}

	want := []EventCategory{
		{
			ID:       "0",
			Owner:    "anodot",
			ImageURL: "https://s3.amazonaws.com/anodot-images-common/logo-anodot.png",
			Name:     "deployments",
		},
		{

			ID:       "1",
			Owner:    "anodot",
			ImageURL: "https://s3.amazonaws.com/anodot-images-common/logo-anodot.png",
			Name:     "alerts",
		},
	}
	if !reflect.DeepEqual(artifacts, want) {
		t.Errorf("eventsService.ListCategories returned %+v, want %+v", artifacts, want)
	}
}

func setup() (c *client.AnodotClient, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	u, _ := url.Parse(server.URL)
	c, err := client.NewAnodotClient(u.String(), "test_1234", nil)
	if err != nil {
		panic(err)
	}

	return c, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testTokenPresent(t *testing.T, r *http.Request, want string) {
	t.Helper()
	values, _ := url.ParseQuery(r.URL.RawQuery)
	actualToken := values.Get("token")
	if actualToken != want {
		t.Errorf("Wrong api token: '%s', want: '%s'", actualToken, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	equal, err := JSONBytesEqual(body, []byte(want))
	if err != nil {
		t.Fatal(err)
	}

	if !equal {
		t.Fatalf("Request body = %+v, want: %+v", string(body), want)
	}
}

func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
