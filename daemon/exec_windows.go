package daemon

import (
	"context"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/rumpl/bof/container"
)

func (daemon *Daemon) execSetPlatformOpt(ctx context.Context, ec *container.ExecConfig, p *specs.Process) error {
	if ec.Container.OS == "windows" {
		p.User.Username = ec.User
	}
	return nil
}
