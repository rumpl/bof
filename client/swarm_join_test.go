package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rumpl/bof/api/types/swarm"
	"github.com/rumpl/bof/errdefs"
)

func TestSwarmJoinError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	err := client.SwarmJoin(context.Background(), swarm.JoinRequest{})
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestSwarmJoin(t *testing.T) {
	expectedURL := "/swarm/join"

	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			if req.Method != http.MethodPost {
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}),
	}

	err := client.SwarmJoin(context.Background(), swarm.JoinRequest{
		ListenAddr: "0.0.0.0:2377",
	})
	if err != nil {
		t.Fatal(err)
	}
}
