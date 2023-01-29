//go:build linux || freebsd || darwin || openbsd
// +build linux freebsd darwin openbsd

package layer

import "github.com/rumpl/bof/pkg/stringid"

func (ls *layerStore) mountID(name string) string {
	return stringid.GenerateRandomID()
}
