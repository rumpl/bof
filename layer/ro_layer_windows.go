package layer // import "github.com/rumpl/bof/layer"

import "github.com/docker/distribution"

var _ distribution.Describable = &roLayer{}

func (rl *roLayer) Descriptor() distribution.Descriptor {
	return rl.descriptor
}
