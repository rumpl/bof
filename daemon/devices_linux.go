package daemon

import (
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/rumpl/bof/api/types/container"
	"github.com/rumpl/bof/pkg/capabilities"
)

var deviceDrivers = map[string]*deviceDriver{}

type deviceDriver struct {
	capset     capabilities.Set
	updateSpec func(*specs.Spec, *deviceInstance) error
}

type deviceInstance struct {
	req          container.DeviceRequest
	selectedCaps []string
}

func registerDeviceDriver(name string, d *deviceDriver) {
	deviceDrivers[name] = d
}

func (daemon *Daemon) handleDevice(req container.DeviceRequest, spec *specs.Spec) error {
	if req.Driver == "" {
		for _, dd := range deviceDrivers {
			if selected := dd.capset.Match(req.Capabilities); selected != nil {
				return dd.updateSpec(spec, &deviceInstance{req: req, selectedCaps: selected})
			}
		}
	} else if dd := deviceDrivers[req.Driver]; dd != nil {
		if selected := dd.capset.Match(req.Capabilities); selected != nil {
			return dd.updateSpec(spec, &deviceInstance{req: req, selectedCaps: selected})
		}
	}
	return incompatibleDeviceRequest{req.Driver, req.Capabilities}
}
