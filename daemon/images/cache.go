package images // import "github.com/rumpl/bof/daemon/images"

import (
	"context"

	"github.com/pkg/errors"
	imagetypes "github.com/rumpl/bof/api/types/image"
	"github.com/rumpl/bof/builder"
	"github.com/rumpl/bof/image/cache"
	"github.com/sirupsen/logrus"
)

// MakeImageCache creates a stateful image cache.
func (i *ImageService) MakeImageCache(ctx context.Context, sourceRefs []string) (builder.ImageCache, error) {
	if len(sourceRefs) == 0 {
		return cache.NewLocal(i.imageStore), nil
	}

	cache := cache.New(i.imageStore)

	for _, ref := range sourceRefs {
		img, err := i.GetImage(ctx, ref, imagetypes.GetImageOpts{})
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			logrus.Warnf("Could not look up %s for cache resolution, skipping: %+v", ref, err)
			continue
		}
		cache.Populate(img)
	}

	return cache, nil
}
