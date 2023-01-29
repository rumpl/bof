package daemon

import (
	"context"

	"github.com/containerd/containerd/containers"
	coci "github.com/containerd/containerd/oci"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/rumpl/bof/container"
)

// WithConsoleSize sets the initial console size
func WithConsoleSize(c *container.Container) coci.SpecOpts {
	return func(ctx context.Context, _ coci.Client, _ *containers.Container, s *coci.Spec) error {
		if c.HostConfig.ConsoleSize[0] > 0 || c.HostConfig.ConsoleSize[1] > 0 {
			s.Process.ConsoleSize = &specs.Box{
				Height: c.HostConfig.ConsoleSize[0],
				Width:  c.HostConfig.ConsoleSize[1],
			}
		}
		return nil
	}
}
