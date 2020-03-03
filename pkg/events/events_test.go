package events

import (
	"encoding/json"
	"fmt"
	"github.com/anodot/anodot-common/pkg/client"
	"github.com/anodot/anodot-common/pkg/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/user-events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testTokenPresent(t, r)

		expectedBody := `{
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		equal, err := JSONBytesEqual(body, []byte(expectedBody))
		if err != nil {
			t.Fatal(err)
		}

		if !equal {
			t.Fatalf("Request body = %+v, want: %+v", string(body), expectedBody)
		}

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

	eventsService := EventsService{c}

	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	ts, err := time.Parse(layout, str)
	if err != nil {
		t.Fatal(err)
	}

	properties := []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{{Key: "service", Value: "myService"}}

	_, err = eventsService.Create(Event{
		Title:       "deployment started on myServer",
		Description: "my description",
		Category:    "deployments",
		Source:      "chef",
		Properties:  properties,
		StartDate:   common.AnodotTimestamp{ts},
		EndDate:     nil,
	})
	if err != nil {
		t.Errorf("eventsService.ListCategories returned error: %v", err)
	}

}

func TestListCategories(t *testing.T) {
	os.Setenv("ANODOT_HTTP_DEBUG_ENABLED", "true")
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/user-events/categories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testTokenPresent(t, r)
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

	eventsService := EventsService{c}

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

	log.Println(u)

	c, err := client.NewAnodotClient(*u, "test", nil)
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

func testTokenPresent(t *testing.T, r *http.Request) {
	t.Helper()
	values, _ := url.ParseQuery(r.URL.RawQuery)
	gotToken := values.Get("token")
	if gotToken != "test" {
		t.Errorf("Request method: %v, want %v", gotToken, "test")
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
