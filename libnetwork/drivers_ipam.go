package libnetwork

import (
	"github.com/rumpl/bof/libnetwork/ipamapi"
	builtinIpam "github.com/rumpl/bof/libnetwork/ipams/builtin"
	nullIpam "github.com/rumpl/bof/libnetwork/ipams/null"
	remoteIpam "github.com/rumpl/bof/libnetwork/ipams/remote"
	"github.com/rumpl/bof/libnetwork/ipamutils"
	"github.com/rumpl/bof/pkg/plugingetter"
)

func initIPAMDrivers(r ipamapi.Registerer, pg plugingetter.PluginGetter, addressPool []*ipamutils.NetworkToSplit) error {
	// TODO: pass address pools as arguments to builtinIpam.Init instead of
	// indirectly through global mutable state. Swarmkit references that
	// function so changing its signature breaks the build.
	if err := builtinIpam.SetDefaultIPAddressPool(addressPool); err != nil {
		return err
	}

	for _, fn := range [](func(ipamapi.Registerer) error){
		builtinIpam.Register,
		nullIpam.Register,
	} {
		if err := fn(r); err != nil {
			return err
		}
	}

	return remoteIpam.Register(r, pg)
}
