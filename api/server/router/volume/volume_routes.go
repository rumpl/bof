package volume // import "github.com/rumpl/bof/api/server/router/volume"

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rumpl/bof/api/server/httputils"
	"github.com/rumpl/bof/api/types/filters"
	"github.com/rumpl/bof/api/types/versions"
	"github.com/rumpl/bof/api/types/volume"
	"github.com/rumpl/bof/volume/service/opts"
	"github.com/sirupsen/logrus"
)

func (v *volumeRouter) getVolumesList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	filters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil {
		return errors.Wrap(err, "error reading volume filters")
	}
	volumes, warnings, err := v.backend.List(ctx, filters)
	if err != nil {
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, &volume.ListResponse{Volumes: volumes, Warnings: warnings})
}

func (v *volumeRouter) getVolumeByName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	vol, err := v.backend.Get(ctx, vars["name"], opts.WithGetResolveStatus)
	if err != nil {
		// otherwise, if this isn't NotFound, or this isn't a high enough version,
		// just return the error by itself.
		return err
	}

	return httputils.WriteJSON(w, http.StatusOK, vol)
}

func (v *volumeRouter) postVolumesCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	var req volume.CreateOptions
	if err := httputils.ReadJSON(r, &req); err != nil {
		return err
	}

	var (
		vol *volume.Volume
		err error
	)
	logrus.Debug("using regular volume")
	vol, err = v.backend.Create(ctx, req.Name, req.Driver, opts.WithCreateOptions(req.DriverOpts), opts.WithCreateLabels(req.Labels))

	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusCreated, vol)
}

func (v *volumeRouter) deleteVolumes(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	force := httputils.BoolValue(r, "force")

	err := v.backend.Remove(ctx, vars["name"], opts.WithPurgeOnError(force))
	if err != nil || force {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (v *volumeRouter) postVolumesPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	pruneFilters, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil {
		return err
	}

	// API version 1.42 changes behavior where prune should only prune anonymous volumes.
	// To keep older API behavior working, we need to add this filter option to consider all (local) volumes for pruning, not just anonymous ones.
	if versions.LessThan(httputils.VersionFromContext(ctx), "1.42") {
		pruneFilters.Add("all", "true")
	}

	pruneReport, err := v.backend.Prune(ctx, pruneFilters)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, pruneReport)
}
