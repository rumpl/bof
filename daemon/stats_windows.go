package daemon

import (
	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/container"
)

// Windows network stats are obtained directly through HCS, hence this is a no-op.
func (daemon *Daemon) getNetworkStats(c *container.Container) (map[string]types.NetworkStats, error) {
	return make(map[string]types.NetworkStats), nil
}
