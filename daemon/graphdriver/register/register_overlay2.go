//go:build !exclude_graphdriver_overlay2 && linux
// +build !exclude_graphdriver_overlay2,linux

package register // import "github.com/rumpl/bof/daemon/graphdriver/register"

import (
	// register the overlay2 graphdriver
	_ "github.com/rumpl/bof/daemon/graphdriver/overlay2"
)
