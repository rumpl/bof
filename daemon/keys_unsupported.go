//go:build !linux
// +build !linux

package daemon // import "github.com/rumpl/bof/daemon"

// modifyRootKeyLimit is a noop on unsupported platforms.
func modifyRootKeyLimit() error {
	return nil
}
