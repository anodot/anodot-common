package anodot_common

import (
	"github.com/anodot/anodot-common/pkg/api"
	"github.com/anodot/anodot-common/pkg/client"
	"log"
	"net/url"
	"os"
	"testing"
)

func TestTse(t *testing.T) {
	os.Setenv("ANODOT_HTTP_DEBUG_ENABLED", "true")
	parse, _ := url.Parse("http://onprem-brian-v3-01.ano-dev.com/")
	anodotClient, err := client.NewAnodotClient(*parse, "a92ee26c2c3373e23e69f437026802ef", nil)
	if err != nil {
		t.Fatal(err)
	}
	parse, _ = url.Parse("https://upload.wikimedia.org/wikipedia/commons/thumb/3/39/Kubernetes_logo_without_workmark.svg/1200px-Kubernetes_logo_without_workmark.svg.png")

	api := api.NewApiClient(anodotClient)

	listSources, err := api.Events.ListSources()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range listSources {
		log.Println(v)
	}
}
