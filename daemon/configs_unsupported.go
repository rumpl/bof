//go:build !linux && !windows
// +build !linux,!windows

package daemon // import "github.com/rumpl/bof/daemon"

func configsSupported() bool {
	return false
}
