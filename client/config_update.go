package client

import (
	"context"
	"net/url"

	"github.com/rumpl/bof/api/types/swarm"
)

// ConfigUpdate attempts to update a config
func (cli *Client) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error {
	if err := cli.NewVersionError("1.30", "config update"); err != nil {
		return err
	}
	query := url.Values{}
	query.Set("version", version.String())
	resp, err := cli.post(ctx, "/configs/"+id+"/update", query, config, nil)
	ensureReaderClosed(resp)
	return err
}
