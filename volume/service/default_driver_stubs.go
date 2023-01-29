//go:build !linux && !windows
// +build !linux,!windows

package service // import "github.com/rumpl/bof/volume/service"

import (
	"github.com/rumpl/bof/pkg/idtools"
	"github.com/rumpl/bof/volume/drivers"
)

func setupDefaultDriver(_ *drivers.Store, _ string, _ idtools.Identity) error { return nil }
