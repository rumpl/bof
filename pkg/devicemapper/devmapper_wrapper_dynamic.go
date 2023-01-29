//go:build linux && cgo && !static_build
// +build linux,cgo,!static_build

package devicemapper // import "github.com/rumpl/bof/pkg/devicemapper"

// #cgo pkg-config: devmapper
import "C"
