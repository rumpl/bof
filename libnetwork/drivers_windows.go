package libnetwork

import (
	"github.com/rumpl/bof/libnetwork/drivers/null"
	"github.com/rumpl/bof/libnetwork/drivers/windows"
	"github.com/rumpl/bof/libnetwork/drivers/windows/overlay"
)

func getInitializers() []initializer {
	return []initializer{
		{null.Register, "null"},
		{overlay.Register, "overlay"},
		{windows.GetInit("transparent"), "transparent"},
		{windows.GetInit("l2bridge"), "l2bridge"},
		{windows.GetInit("l2tunnel"), "l2tunnel"},
		{windows.GetInit("nat"), "nat"},
		{windows.GetInit("internal"), "internal"},
		{windows.GetInit("private"), "private"},
		{windows.GetInit("ics"), "ics"},
	}
}
