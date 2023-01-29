package daemon // import "github.com/rumpl/bof/daemon"

import (
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/rumpl/bof/container"
)

func setLinuxDomainname(c *container.Container, s *specs.Spec) {
	// There isn't a field in the OCI for the NIS domainname, but luckily there
	// is a sysctl which has an identical effect to setdomainname(2) so there's
	// no explicit need for runtime support.
	s.Linux.Sysctl = make(map[string]string)
	if c.Config.Domainname != "" {
		s.Linux.Sysctl["kernel.domainname"] = c.Config.Domainname
	}
}
