//go:build !exclude_graphdriver_overlay && linux
// +build !exclude_graphdriver_overlay,linux

package register

import (
	// register the overlay graphdriver
	_ "github.com/rumpl/bof/daemon/graphdriver/overlay"
)
