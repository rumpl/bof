package client // import "github.com/rumpl/bof/client"

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/filters"
	"github.com/rumpl/bof/api/types/registry"
	"github.com/rumpl/bof/errdefs"
)

// ImageSearch makes the docker host search by a term in a remote registry.
// The list of results is not sorted in any fashion.
func (cli *Client) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	var results []registry.SearchResult
	query := url.Values{}
	query.Set("term", term)
	if options.Limit > 0 {
		query.Set("limit", strconv.Itoa(options.Limit))
	}

	if options.Filters.Len() > 0 {
		filterJSON, err := filters.ToJSON(options.Filters)
		if err != nil {
			return results, err
		}
		query.Set("filters", filterJSON)
	}

	resp, err := cli.tryImageSearch(ctx, query, options.RegistryAuth)
	defer ensureReaderClosed(resp)
	if errdefs.IsUnauthorized(err) && options.PrivilegeFunc != nil {
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil {
			return results, privilegeErr
		}
		resp, err = cli.tryImageSearch(ctx, query, newAuthHeader)
	}
	if err != nil {
		return results, err
	}

	err = json.NewDecoder(resp.body).Decode(&results)
	return results, err
}

func (cli *Client) tryImageSearch(ctx context.Context, query url.Values, registryAuth string) (serverResponse, error) {
	headers := map[string][]string{registry.AuthHeader: {registryAuth}}
	return cli.get(ctx, "/images/search", query, headers)
}
