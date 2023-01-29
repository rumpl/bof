package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/fs"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	containertypes "github.com/rumpl/bof/api/types/container"
	"github.com/rumpl/bof/container"
	"golang.org/x/sys/windows/registry"
	"gotest.tools/v3/assert"
)

func TestSetupWindowsDevices(t *testing.T) {
	t.Run("it does nothing if there are no devices", func(t *testing.T) {
		devices, err := setupWindowsDevices(nil)
		assert.NilError(t, err)
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it fails if any devices are blank", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class/anything"}, {PathOnHost: ""}})
		assert.ErrorContains(t, err, "invalid device assignment path")
		assert.ErrorContains(t, err, "''")
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it fails if all devices do not contain '/' or '://'", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "anything"}, {PathOnHost: "goes"}})
		assert.ErrorContains(t, err, "invalid device assignment path")
		assert.ErrorContains(t, err, "'anything'")
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it fails if any devices do not contain '/' or '://'", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class/anything"}, {PathOnHost: "goes"}})
		assert.ErrorContains(t, err, "invalid device assignment path")
		assert.ErrorContains(t, err, "'goes'")
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it fails if all '/'-separated devices do not have IDType 'class'", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "klass/anything"}, {PathOnHost: "klass/goes"}})
		assert.ErrorContains(t, err, "invalid device assignment path")
		assert.ErrorContains(t, err, "'klass/anything'")
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it fails if any '/'-separated devices do not have IDType 'class'", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class/anything"}, {PathOnHost: "klass/goes"}})
		assert.ErrorContains(t, err, "invalid device assignment path")
		assert.ErrorContains(t, err, "'klass/goes'")
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it fails if any '://'-separated devices have IDType ''", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class/anything"}, {PathOnHost: "://goes"}})
		assert.ErrorContains(t, err, "invalid device assignment path")
		assert.ErrorContains(t, err, "'://goes'")
		assert.Equal(t, len(devices), 0)
	})

	t.Run("it creates devices if all '/'-separated devices have IDType 'class'", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class/anything"}, {PathOnHost: "class/goes"}})
		expectedDevices := []specs.WindowsDevice{{IDType: "class", ID: "anything"}, {IDType: "class", ID: "goes"}}
		assert.NilError(t, err)
		assert.Equal(t, len(devices), len(expectedDevices))
		for i := range expectedDevices {
			assert.Equal(t, devices[i], expectedDevices[i])
		}
	})

	t.Run("it creates devices if all '://'-separated devices have non-blank IDType", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class://anything"}, {PathOnHost: "klass://goes"}})
		expectedDevices := []specs.WindowsDevice{{IDType: "class", ID: "anything"}, {IDType: "klass", ID: "goes"}}
		assert.NilError(t, err)
		assert.Equal(t, len(devices), len(expectedDevices))
		for i := range expectedDevices {
			assert.Equal(t, devices[i], expectedDevices[i])
		}
	})

	t.Run("it creates devices when given a mix of '/'-separated and '://'-separated devices", func(t *testing.T) {
		devices, err := setupWindowsDevices([]containertypes.DeviceMapping{{PathOnHost: "class/anything"}, {PathOnHost: "klass://goes"}})
		expectedDevices := []specs.WindowsDevice{{IDType: "class", ID: "anything"}, {IDType: "klass", ID: "goes"}}
		assert.NilError(t, err)
		assert.Equal(t, len(devices), len(expectedDevices))
		for i := range expectedDevices {
			assert.Equal(t, devices[i], expectedDevices[i])
		}
	})
}
