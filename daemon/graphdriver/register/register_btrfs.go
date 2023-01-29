//go:build !exclude_graphdriver_btrfs && linux
// +build !exclude_graphdriver_btrfs,linux

package register // import "github.com/rumpl/bof/daemon/graphdriver/register"

import (
	// register the btrfs graphdriver
	_ "github.com/rumpl/bof/daemon/graphdriver/btrfs"
)
