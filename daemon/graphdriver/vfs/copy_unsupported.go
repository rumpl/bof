//go:build !linux
// +build !linux

package vfs // import "github.com/rumpl/bof/daemon/graphdriver/vfs"

import (
	"github.com/rumpl/bof/pkg/chrootarchive"
	"github.com/rumpl/bof/pkg/idtools"
)

func dirCopy(srcDir, dstDir string) error {
	return chrootarchive.NewArchiver(idtools.IdentityMapping{}).CopyWithTar(srcDir, dstDir)
}
