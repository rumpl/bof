//go:build !linux
// +build !linux

package caps

func initCaps() {
	// no capabilities on Windows
}
