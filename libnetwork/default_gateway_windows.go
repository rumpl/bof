package libnetwork

import (
	windriver "github.com/rumpl/bof/libnetwork/drivers/windows"
	"github.com/rumpl/bof/libnetwork/options"
	"github.com/rumpl/bof/libnetwork/types"
)

const libnGWNetwork = "nat"

func getPlatformOption() EndpointOption {

	epOption := options.Generic{
		windriver.DisableICC: true,
		windriver.DisableDNS: true,
	}
	return EndpointOptionGeneric(epOption)
}

func (c *Controller) createGWNetwork() (Network, error) {
	return nil, types.NotImplementedErrorf("default gateway functionality is not implemented in windows")
}
