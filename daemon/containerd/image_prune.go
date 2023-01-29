package containerd

import (
	"context"
	"errors"

	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/filters"
	"github.com/rumpl/bof/errdefs"
)

// ImagesPrune removes unused images
func (i *ImageService) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (*types.ImagesPruneReport, error) {
	return nil, errdefs.NotImplemented(errors.New("not implemented"))
}
