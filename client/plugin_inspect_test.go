package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/errdefs"
)

func TestPluginInspectError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	_, _, err := client.PluginInspectWithRaw(context.Background(), "nothing")
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestPluginInspectWithEmptyID(t *testing.T) {
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("should not make request")
		}),
	}
	_, _, err := client.PluginInspectWithRaw(context.Background(), "")
	if !IsErrNotFound(err) {
		t.Fatalf("Expected NotFoundError, got %v", err)
	}
}

func TestPluginInspect(t *testing.T) {
	expectedURL := "/plugins/plugin_name"
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			content, err := json.Marshal(types.Plugin{
				ID: "plugin_id",
			})
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(content)),
			}, nil
		}),
	}

	pluginInspect, _, err := client.PluginInspectWithRaw(context.Background(), "plugin_name")
	if err != nil {
		t.Fatal(err)
	}
	if pluginInspect.ID != "plugin_id" {
		t.Fatalf("expected `plugin_id`, got %s", pluginInspect.ID)
	}
}
