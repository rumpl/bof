package client // import "github.com/rumpl/bof/client"

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rumpl/bof/api/types"
)

// Info returns information about the docker server.
func (cli *Client) Info(ctx context.Context) (types.Info, error) {
	var info types.Info
	serverResp, err := cli.get(ctx, "/info", url.Values{}, nil)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return info, err
	}

	if err := json.NewDecoder(serverResp.body).Decode(&info); err != nil {
		return info, fmt.Errorf("Error reading remote info: %v", err)
	}

	return info, nil
}
