package daemon

import (
	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/dockerversion"
)

func (daemon *Daemon) fillLicense(v *types.Info) {
	v.ProductLicense = dockerversion.DefaultProductLicense
}
