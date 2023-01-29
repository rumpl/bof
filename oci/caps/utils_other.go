//go:build !linux
// +build !linux

package caps // import "github.com/rumpl/bof/oci/caps"

func initCaps() {
	// no capabilities on Windows
}
