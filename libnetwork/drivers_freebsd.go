package libnetwork

import (
	"github.com/rumpl/bof/libnetwork/drivers/null"
)

func getInitializers() []initializer {
	return []initializer{
		{null.Register, "null"},
	}
}
