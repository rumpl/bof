package daemon

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/events"
	"github.com/rumpl/bof/api/types/filters"
	timetypes "github.com/rumpl/bof/api/types/time"
	"github.com/rumpl/bof/errdefs"
	"github.com/rumpl/bof/libnetwork"
	"github.com/rumpl/bof/runconfig"
	"github.com/sirupsen/logrus"
)

var (
	// errPruneRunning is returned when a prune request is received while
	// one is in progress
	errPruneRunning = errdefs.Conflict(errors.New("a prune operation is already running"))

	containersAcceptedFilters = map[string]bool{
		"label":  true,
		"label!": true,
		"until":  true,
	}

	networksAcceptedFilters = map[string]bool{
		"label":  true,
		"label!": true,
		"until":  true,
	}
)

// ContainersPrune removes unused containers
func (daemon *Daemon) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (*types.ContainersPruneReport, error) {
	if !atomic.CompareAndSwapInt32(&daemon.pruneRunning, 0, 1) {
		return nil, errPruneRunning
	}
	defer atomic.StoreInt32(&daemon.pruneRunning, 0)

	rep := &types.ContainersPruneReport{}

	// make sure that only accepted filters have been received
	err := pruneFilters.Validate(containersAcceptedFilters)
	if err != nil {
		return nil, err
	}

	until, err := getUntilFromPruneFilters(pruneFilters)
	if err != nil {
		return nil, err
	}

	allContainers := daemon.List()
	for _, c := range allContainers {
		select {
		case <-ctx.Done():
			logrus.Debugf("ContainersPrune operation cancelled: %#v", *rep)
			return rep, nil
		default:
		}

		if !c.IsRunning() {
			if !until.IsZero() && c.Created.After(until) {
				continue
			}
			if !matchLabels(pruneFilters, c.Config.Labels) {
				continue
			}
			cSize, _ := daemon.imageService.GetContainerLayerSize(c.ID)
			// TODO: sets RmLink to true?
			err := daemon.ContainerRm(c.ID, &types.ContainerRmConfig{})
			if err != nil {
				logrus.Warnf("failed to prune container %s: %v", c.ID, err)
				continue
			}
			if cSize > 0 {
				rep.SpaceReclaimed += uint64(cSize)
			}
			rep.ContainersDeleted = append(rep.ContainersDeleted, c.ID)
		}
	}
	daemon.EventsService.Log("prune", events.ContainerEventType, events.Actor{
		Attributes: map[string]string{"reclaimed": strconv.FormatUint(rep.SpaceReclaimed, 10)},
	})
	return rep, nil
}

// localNetworksPrune removes unused local networks
func (daemon *Daemon) localNetworksPrune(ctx context.Context, pruneFilters filters.Args) *types.NetworksPruneReport {
	rep := &types.NetworksPruneReport{}

	until, _ := getUntilFromPruneFilters(pruneFilters)

	// When the function returns true, the walk will stop.
	l := func(nw libnetwork.Network) bool {
		select {
		case <-ctx.Done():
			// context cancelled
			return true
		default:
		}
		if nw.Info().ConfigOnly() {
			return false
		}
		if !until.IsZero() && nw.Info().Created().After(until) {
			return false
		}
		if !matchLabels(pruneFilters, nw.Info().Labels()) {
			return false
		}
		nwName := nw.Name()
		if runconfig.IsPreDefinedNetwork(nwName) {
			return false
		}
		if len(nw.Endpoints()) > 0 {
			return false
		}
		if err := daemon.DeleteNetwork(nw.ID()); err != nil {
			logrus.Warnf("could not remove local network %s: %v", nwName, err)
			return false
		}
		rep.NetworksDeleted = append(rep.NetworksDeleted, nwName)
		return false
	}
	daemon.netController.WalkNetworks(l)
	return rep
}

// NetworksPrune removes unused networks
func (daemon *Daemon) NetworksPrune(ctx context.Context, pruneFilters filters.Args) (*types.NetworksPruneReport, error) {
	if !atomic.CompareAndSwapInt32(&daemon.pruneRunning, 0, 1) {
		return nil, errPruneRunning
	}
	defer atomic.StoreInt32(&daemon.pruneRunning, 0)

	// make sure that only accepted filters have been received
	err := pruneFilters.Validate(networksAcceptedFilters)
	if err != nil {
		return nil, err
	}

	if _, err := getUntilFromPruneFilters(pruneFilters); err != nil {
		return nil, err
	}

	rep := &types.NetworksPruneReport{}

	localRep := daemon.localNetworksPrune(ctx, pruneFilters)
	rep.NetworksDeleted = append(rep.NetworksDeleted, localRep.NetworksDeleted...)

	select {
	case <-ctx.Done():
		logrus.Debugf("NetworksPrune operation cancelled: %#v", *rep)
		return rep, nil
	default:
	}
	daemon.EventsService.Log("prune", events.NetworkEventType, events.Actor{
		Attributes: map[string]string{"reclaimed": "0"},
	})
	return rep, nil
}

func getUntilFromPruneFilters(pruneFilters filters.Args) (time.Time, error) {
	until := time.Time{}
	if !pruneFilters.Contains("until") {
		return until, nil
	}
	untilFilters := pruneFilters.Get("until")
	if len(untilFilters) > 1 {
		return until, fmt.Errorf("more than one until filter specified")
	}
	ts, err := timetypes.GetTimestamp(untilFilters[0], time.Now())
	if err != nil {
		return until, err
	}
	seconds, nanoseconds, err := timetypes.ParseTimestamps(ts, 0)
	if err != nil {
		return until, err
	}
	until = time.Unix(seconds, nanoseconds)
	return until, nil
}

func matchLabels(pruneFilters filters.Args, labels map[string]string) bool {
	if !pruneFilters.MatchKVList("label", labels) {
		return false
	}
	// By default MatchKVList will return true if field (like 'label!') does not exist
	// So we have to add additional Contains("label!") check
	if pruneFilters.Contains("label!") {
		if pruneFilters.MatchKVList("label!", labels) {
			return false
		}
	}
	return true
}
