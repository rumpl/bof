package libnetwork

import (
	"github.com/rumpl/bof/libnetwork/drivers/bridge"
	"github.com/rumpl/bof/libnetwork/drivers/host"
	"github.com/rumpl/bof/libnetwork/drivers/ipvlan"
	"github.com/rumpl/bof/libnetwork/drivers/macvlan"
	"github.com/rumpl/bof/libnetwork/drivers/null"
	"github.com/rumpl/bof/libnetwork/drivers/overlay"
)

func getInitializers() []initializer {
	in := []initializer{
		{bridge.Register, "bridge"},
		{host.Register, "host"},
		{ipvlan.Register, "ipvlan"},
		{macvlan.Register, "macvlan"},
		{null.Register, "null"},
		{overlay.Register, "overlay"},
	}
	return in
}
