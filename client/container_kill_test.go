package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rumpl/bof/errdefs"
)

func TestContainerKillError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}
	err := client.ContainerKill(context.Background(), "nothing", "SIGKILL")
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestContainerKill(t *testing.T) {
	expectedURL := "/containers/container_id/kill"
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			signal := req.URL.Query().Get("signal")
			if signal != "SIGKILL" {
				return nil, fmt.Errorf("signal not set in URL query properly. Expected 'SIGKILL', got %s", signal)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}),
	}

	err := client.ContainerKill(context.Background(), "container_id", "SIGKILL")
	if err != nil {
		t.Fatal(err)
	}
}
