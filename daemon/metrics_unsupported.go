//go:build windows
// +build windows

package daemon // import "github.com/rumpl/bof/daemon"

import "github.com/rumpl/bof/pkg/plugingetter"

func registerMetricsPluginCallback(getter plugingetter.PluginGetter, sockPath string) {
}

func (daemon *Daemon) listenMetricsSock() (string, error) {
	return "", nil
}
