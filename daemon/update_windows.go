package daemon

import (
	"github.com/rumpl/bof/api/types/container"
	libcontainerdtypes "github.com/rumpl/bof/libcontainerd/types"
)

func toContainerdResources(resources container.Resources) *libcontainerdtypes.Resources {
	// We don't support update, so do nothing
	return nil
}
