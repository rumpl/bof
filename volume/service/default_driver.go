//go:build linux || windows

package service

import (
	"github.com/pkg/errors"
	"github.com/rumpl/bof/pkg/idtools"
	"github.com/rumpl/bof/volume"
	"github.com/rumpl/bof/volume/drivers"
	"github.com/rumpl/bof/volume/local"
)

func setupDefaultDriver(store *drivers.Store, root string, rootIDs idtools.Identity) error {
	d, err := local.New(root, rootIDs)
	if err != nil {
		return errors.Wrap(err, "error setting up default driver")
	}
	if !store.Register(d, volume.DefaultDriverName) {
		return errors.New("local volume driver could not be registered")
	}
	return nil
}
