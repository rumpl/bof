//go:build !exclude_graphdriver_fuseoverlayfs && linux
// +build !exclude_graphdriver_fuseoverlayfs,linux

package register // import "github.com/rumpl/bof/daemon/graphdriver/register"

import (
	// register the fuse-overlayfs graphdriver
	_ "github.com/rumpl/bof/daemon/graphdriver/fuse-overlayfs"
)
