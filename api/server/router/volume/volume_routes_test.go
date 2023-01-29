package volume

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/filters"
	"github.com/rumpl/bof/api/types/volume"
	"github.com/rumpl/bof/errdefs"
	"github.com/rumpl/bof/volume/service/opts"
)

func callGetVolume(v *volumeRouter, name string) (*httptest.ResponseRecorder, error) {
	vars := map[string]string{"name": name}
	req := httptest.NewRequest("GET", fmt.Sprintf("/volumes/%s", name), nil)
	resp := httptest.NewRecorder()

	err := v.getVolumeByName(context.Background(), resp, req, vars)
	return resp, err
}

func callListVolumes(v *volumeRouter) (*httptest.ResponseRecorder, error) {
	vars := map[string]string{}
	req := httptest.NewRequest("GET", "/volumes", nil)
	resp := httptest.NewRecorder()

	err := v.getVolumesList(context.Background(), resp, req, vars)
	return resp, err
}

func TestGetVolumeByNameNotFoundNoSwarm(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{},
	}

	_, err := callGetVolume(v, "notReal")

	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsNotFound(err))
}

func TestGetVolumeByNameNotFoundNotManager(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{},
	}

	_, err := callGetVolume(v, "notReal")

	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsNotFound(err))
}

func TestGetVolumeByNameNotFound(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{},
	}

	_, err := callGetVolume(v, "notReal")

	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsNotFound(err))
}

func TestGetVolumeByNameFoundRegular(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{
			volumes: map[string]*volume.Volume{

				"volume1": {
					Name: "volume1",
				},
			},
		},
	}

	_, err := callGetVolume(v, "volume1")
	assert.NilError(t, err)
}

func TestGetVolumeByNameFoundSwarm(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{},
	}

	_, err := callGetVolume(v, "volume1")
	assert.NilError(t, err)
}
func TestListVolumes(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{
			volumes: map[string]*volume.Volume{
				"v1": {Name: "v1"},
				"v2": {Name: "v2"},
			},
		},
	}

	resp, err := callListVolumes(v)
	assert.NilError(t, err)
	d := json.NewDecoder(resp.Result().Body)
	respVols := volume.ListResponse{}
	assert.NilError(t, d.Decode(&respVols))

	assert.Assert(t, respVols.Volumes != nil)
	assert.Equal(t, len(respVols.Volumes), 4, "volumes %v", respVols.Volumes)
}

func TestListVolumesNoSwarm(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{
			volumes: map[string]*volume.Volume{
				"v1": {Name: "v1"},
				"v2": {Name: "v2"},
			},
		},
	}

	_, err := callListVolumes(v)
	assert.NilError(t, err)
}

func TestListVolumesNoManager(t *testing.T) {
	v := &volumeRouter{
		backend: &fakeVolumeBackend{
			volumes: map[string]*volume.Volume{
				"v1": {Name: "v1"},
				"v2": {Name: "v2"},
			},
		},
	}

	resp, err := callListVolumes(v)
	assert.NilError(t, err)

	d := json.NewDecoder(resp.Result().Body)
	respVols := volume.ListResponse{}
	assert.NilError(t, d.Decode(&respVols))

	assert.Equal(t, len(respVols.Volumes), 2)
	assert.Equal(t, len(respVols.Warnings), 0)
}

func TestCreateRegularVolume(t *testing.T) {
	b := &fakeVolumeBackend{}
	c := &fakeClusterBackend{
		swarm:   true,
		manager: true,
	}
	v := &volumeRouter{
		backend: b,
	}

	volumeCreate := volume.CreateOptions{
		Name:   "vol1",
		Driver: "foodriver",
	}

	buf := bytes.Buffer{}
	e := json.NewEncoder(&buf)
	e.Encode(volumeCreate)

	req := httptest.NewRequest("POST", "/volumes/create", &buf)
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	err := v.postVolumesCreate(context.Background(), resp, req, nil)

	assert.NilError(t, err)

	respVolume := volume.Volume{}

	assert.NilError(t, json.NewDecoder(resp.Result().Body).Decode(&respVolume))

	assert.Equal(t, respVolume.Name, "vol1")
	assert.Equal(t, respVolume.Driver, "foodriver")

	assert.Equal(t, 1, len(b.volumes))
	assert.Equal(t, 0, len(c.volumes))
}

func TestCreateSwarmVolumeNoSwarm(t *testing.T) {
	b := &fakeVolumeBackend{}

	v := &volumeRouter{
		backend: b,
	}

	volumeCreate := volume.CreateOptions{
		ClusterVolumeSpec: &volume.ClusterVolumeSpec{},
		Name:              "volCluster",
		Driver:            "someCSI",
	}

	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(volumeCreate)

	req := httptest.NewRequest("POST", "/volumes/create", &buf)
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	err := v.postVolumesCreate(context.Background(), resp, req, nil)

	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsUnavailable(err))
}

func TestCreateSwarmVolumeNotManager(t *testing.T) {
	b := &fakeVolumeBackend{}
	v := &volumeRouter{
		backend: b,
	}

	volumeCreate := volume.CreateOptions{
		ClusterVolumeSpec: &volume.ClusterVolumeSpec{},
		Name:              "volCluster",
		Driver:            "someCSI",
	}

	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(volumeCreate)

	req := httptest.NewRequest("POST", "/volumes/create", &buf)
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	err := v.postVolumesCreate(context.Background(), resp, req, nil)

	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsUnavailable(err))
}

func TestCreateVolumeCluster(t *testing.T) {
	b := &fakeVolumeBackend{}
	v := &volumeRouter{
		backend: b,
	}

	volumeCreate := volume.CreateOptions{
		ClusterVolumeSpec: &volume.ClusterVolumeSpec{},
		Name:              "volCluster",
		Driver:            "someCSI",
	}

	buf := bytes.Buffer{}
	json.NewEncoder(&buf).Encode(volumeCreate)

	req := httptest.NewRequest("POST", "/volumes/create", &buf)
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	err := v.postVolumesCreate(context.Background(), resp, req, nil)

	assert.NilError(t, err)

	respVolume := volume.Volume{}

	assert.NilError(t, json.NewDecoder(resp.Result().Body).Decode(&respVolume))

	assert.Equal(t, respVolume.Name, "volCluster")
	assert.Equal(t, respVolume.Driver, "someCSI")

	assert.Equal(t, 0, len(b.volumes))
}

func TestVolumeRemove(t *testing.T) {
	b := &fakeVolumeBackend{
		volumes: map[string]*volume.Volume{
			"vol1": {
				Name: "vol1",
			},
		},
	}

	v := &volumeRouter{
		backend: b,
	}

	req := httptest.NewRequest("DELETE", "/volumes/vol1", nil)
	resp := httptest.NewRecorder()

	err := v.deleteVolumes(context.Background(), resp, req, map[string]string{"name": "vol1"})
	assert.NilError(t, err)
	assert.Equal(t, len(b.volumes), 0)
}

func TestVolumeRemoveNotFoundNoSwarm(t *testing.T) {
	b := &fakeVolumeBackend{}
	v := &volumeRouter{
		backend: b,
	}

	req := httptest.NewRequest("DELETE", "/volumes/vol1", nil)
	resp := httptest.NewRecorder()

	err := v.deleteVolumes(context.Background(), resp, req, map[string]string{"name": "vol1"})
	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsNotFound(err), err.Error())
}

func TestVolumeRemoveNotFoundNoManager(t *testing.T) {
	b := &fakeVolumeBackend{}
	v := &volumeRouter{
		backend: b,
	}

	req := httptest.NewRequest("DELETE", "/volumes/vol1", nil)
	resp := httptest.NewRecorder()

	err := v.deleteVolumes(context.Background(), resp, req, map[string]string{"name": "vol1"})
	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsNotFound(err))
}

func TestVolumeRemoveFoundNoSwarm(t *testing.T) {
	b := &fakeVolumeBackend{
		volumes: map[string]*volume.Volume{
			"vol1": {
				Name: "vol1",
			},
		},
	}

	v := &volumeRouter{
		backend: b,
	}

	req := httptest.NewRequest("DELETE", "/volumes/vol1", nil)
	resp := httptest.NewRecorder()

	err := v.deleteVolumes(context.Background(), resp, req, map[string]string{"name": "vol1"})
	assert.NilError(t, err)
	assert.Equal(t, len(b.volumes), 0)
}

func TestVolumeRemoveNoSwarmInUse(t *testing.T) {
	b := &fakeVolumeBackend{
		volumes: map[string]*volume.Volume{
			"inuse": {
				Name: "inuse",
			},
		},
	}
	v := &volumeRouter{
		backend: b,
	}

	req := httptest.NewRequest("DELETE", "/volumes/inuse", nil)
	resp := httptest.NewRecorder()

	err := v.deleteVolumes(context.Background(), resp, req, map[string]string{"name": "inuse"})
	assert.Assert(t, err != nil)
	assert.Assert(t, errdefs.IsConflict(err))
}

type fakeVolumeBackend struct {
	volumes map[string]*volume.Volume
}

func (b *fakeVolumeBackend) List(_ context.Context, _ filters.Args) ([]*volume.Volume, []string, error) {
	volumes := []*volume.Volume{}
	for _, v := range b.volumes {
		volumes = append(volumes, v)
	}
	return volumes, nil, nil
}

func (b *fakeVolumeBackend) Get(_ context.Context, name string, _ ...opts.GetOption) (*volume.Volume, error) {
	if v, ok := b.volumes[name]; ok {
		return v, nil
	}
	return nil, errdefs.NotFound(fmt.Errorf("volume %s not found", name))
}

func (b *fakeVolumeBackend) Create(_ context.Context, name, driverName string, _ ...opts.CreateOption) (*volume.Volume, error) {
	if _, ok := b.volumes[name]; ok {
		// TODO(dperny): return appropriate error type
		return nil, fmt.Errorf("already exists")
	}

	v := &volume.Volume{
		Name:   name,
		Driver: driverName,
	}
	if b.volumes == nil {
		b.volumes = map[string]*volume.Volume{
			name: v,
		}
	} else {
		b.volumes[name] = v
	}

	return v, nil
}

func (b *fakeVolumeBackend) Remove(_ context.Context, name string, o ...opts.RemoveOption) error {
	removeOpts := &opts.RemoveConfig{}
	for _, opt := range o {
		opt(removeOpts)
	}

	if v, ok := b.volumes[name]; !ok {
		if !removeOpts.PurgeOnError {
			return errdefs.NotFound(fmt.Errorf("volume %s not found", name))
		}
	} else if v.Name == "inuse" {
		return errdefs.Conflict(fmt.Errorf("volume in use"))
	}

	delete(b.volumes, name)

	return nil
}

func (b *fakeVolumeBackend) Prune(_ context.Context, _ filters.Args) (*types.VolumesPruneReport, error) {
	return nil, nil
}

type fakeClusterBackend struct {
	swarm   bool
	manager bool
	idCount int
	volumes map[string]*volume.Volume
}

func (c *fakeClusterBackend) checkSwarm() error {
	if !c.swarm {
		return errdefs.Unavailable(fmt.Errorf("this node is not a swarm manager. Use \"docker swarm init\" or \"docker swarm join\" to connect this node to swarm and try again"))
	} else if !c.manager {
		return errdefs.Unavailable(fmt.Errorf("this node is not a swarm manager. Worker nodes can't be used to view or modify cluster state. Please run this command on a manager node or promote the current node to a manager"))
	}

	return nil
}

func (c *fakeClusterBackend) IsManager() bool {
	return c.swarm && c.manager
}

func (c *fakeClusterBackend) GetVolume(nameOrID string) (volume.Volume, error) {
	if err := c.checkSwarm(); err != nil {
		return volume.Volume{}, err
	}

	if v, ok := c.volumes[nameOrID]; ok {
		return *v, nil
	}
	return volume.Volume{}, errdefs.NotFound(fmt.Errorf("volume %s not found", nameOrID))
}

func (c *fakeClusterBackend) GetVolumes(options volume.ListOptions) ([]*volume.Volume, error) {
	if err := c.checkSwarm(); err != nil {
		return nil, err
	}

	volumes := []*volume.Volume{}

	for _, v := range c.volumes {
		volumes = append(volumes, v)
	}
	return volumes, nil
}

func (c *fakeClusterBackend) CreateVolume(volumeCreate volume.CreateOptions) (*volume.Volume, error) {
	if err := c.checkSwarm(); err != nil {
		return nil, err
	}

	if _, ok := c.volumes[volumeCreate.Name]; ok {
		// TODO(dperny): return appropriate already exists error
		return nil, fmt.Errorf("already exists")
	}

	v := &volume.Volume{
		Name:    volumeCreate.Name,
		Driver:  volumeCreate.Driver,
		Labels:  volumeCreate.Labels,
		Options: volumeCreate.DriverOpts,
		Scope:   "global",
	}

	c.idCount = c.idCount + 1
	if c.volumes == nil {
		c.volumes = map[string]*volume.Volume{
			v.Name: v,
		}
	} else {
		c.volumes[v.Name] = v
	}

	return v, nil
}

func (c *fakeClusterBackend) RemoveVolume(nameOrID string, force bool) error {
	if err := c.checkSwarm(); err != nil {
		return err
	}

	v, ok := c.volumes[nameOrID]
	if !ok {
		return errdefs.NotFound(fmt.Errorf("volume %s not found", nameOrID))
	}

	if _, mustforce := v.Options["mustforce"]; mustforce && !force {
		return errdefs.Conflict(fmt.Errorf("volume %s must be force removed", nameOrID))
	}

	delete(c.volumes, nameOrID)

	return nil
}

func (c *fakeClusterBackend) UpdateVolume(nameOrID string, version uint64, _ volume.UpdateOptions) error {
	if err := c.checkSwarm(); err != nil {
		return err
	}

	if _, ok := c.volumes[nameOrID]; !ok {

		return errdefs.NotFound(fmt.Errorf("volume %q not found", nameOrID))
	}

	return nil
}
