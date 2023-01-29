//go:build !windows
// +build !windows

package dockerfile // import "github.com/rumpl/bof/builder/dockerfile"

func defaultShellForOS(os string) []string {
	return []string{"/bin/sh", "-c"}
}
