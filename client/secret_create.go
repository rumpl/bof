package client // import "github.com/rumpl/bof/client"

import (
	"context"
	"encoding/json"

	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/swarm"
)

// SecretCreate creates a new secret.
func (cli *Client) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (types.SecretCreateResponse, error) {
	var response types.SecretCreateResponse
	if err := cli.NewVersionError("1.25", "secret create"); err != nil {
		return response, err
	}
	resp, err := cli.post(ctx, "/secrets/create", nil, secret, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	return response, err
}
