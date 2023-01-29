package daemon // import "github.com/rumpl/bof/daemon"

import (
	"context"

	"github.com/rumpl/bof/api/types/registry"
	"github.com/rumpl/bof/dockerversion"
)

// AuthenticateToRegistry checks the validity of credentials in authConfig
func (daemon *Daemon) AuthenticateToRegistry(ctx context.Context, authConfig *registry.AuthConfig) (string, string, error) {
	return daemon.registryService.Auth(ctx, authConfig, dockerversion.DockerUserAgent(ctx))
}
