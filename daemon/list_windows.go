package daemon // import "github.com/rumpl/bof/daemon"

import (
	"strings"

	"github.com/rumpl/bof/container"
)

// excludeByIsolation is a platform specific helper function to support PS
// filtering by Isolation. This is a Windows-only concept, so is a no-op on Unix.
func excludeByIsolation(container *container.Snapshot, ctx *listContext) iterationAction {
	i := strings.ToLower(container.HostConfig.Isolation)
	if i == "" {
		i = "default"
	}
	if !ctx.filters.Match("isolation", i) {
		return excludeContainer
	}
	return includeContainer
}
