package system

import "github.com/rumpl/bof/pkg/sysinfo"

// MemInfo contains memory statistics of the host system.
//
// Deprecated: use [sysinfo.Memory].
type MemInfo = sysinfo.Memory

// ReadMemInfo retrieves memory statistics of the host system and returns a
// MemInfo type.
//
// Deprecated: use [sysinfo.ReadMemInfo].
func ReadMemInfo() (*sysinfo.Memory, error) {
	return sysinfo.ReadMemInfo()
}
