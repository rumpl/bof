package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/filters"
)

// VolumesPrune requests the daemon to delete unused data
func (cli *Client) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (types.VolumesPruneReport, error) {
	var report types.VolumesPruneReport

	if err := cli.NewVersionError("1.25", "volume prune"); err != nil {
		return report, err
	}

	query, err := getFiltersQuery(pruneFilters)
	if err != nil {
		return report, err
	}

	serverResp, err := cli.post(ctx, "/volumes/prune", query, nil, nil)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return report, err
	}

	if err := json.NewDecoder(serverResp.body).Decode(&report); err != nil {
		return report, fmt.Errorf("error retrieving volume prune report: %v", err)
	}

	return report, nil
}
