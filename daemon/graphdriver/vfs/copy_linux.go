package vfs // import "github.com/rumpl/bof/daemon/graphdriver/vfs"

import "github.com/rumpl/bof/daemon/graphdriver/copy"

func dirCopy(srcDir, dstDir string) error {
	return copy.DirCopy(srcDir, dstDir, copy.Content, false)
}
