package containerd

import (
	"context"

	"github.com/rumpl/bof/builder"
)

// MakeImageCache creates a stateful image cache.
func (i *ImageService) MakeImageCache(ctx context.Context, cacheFrom []string) (builder.ImageCache, error) {
	panic("not implemented")
}
