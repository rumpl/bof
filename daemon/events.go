package daemon

import (
	"strings"
	"time"

	gogotypes "github.com/gogo/protobuf/types"
	swarmapi "github.com/moby/swarmkit/v2/api"
	"github.com/rumpl/bof/api/types/events"
	"github.com/rumpl/bof/api/types/filters"
	"github.com/rumpl/bof/container"
	daemonevents "github.com/rumpl/bof/daemon/events"
	"github.com/rumpl/bof/libnetwork"
)

var (
	clusterEventAction = map[swarmapi.WatchActionKind]string{
		swarmapi.WatchActionKindCreate: "create",
		swarmapi.WatchActionKindUpdate: "update",
		swarmapi.WatchActionKindRemove: "remove",
	}
)

// LogContainerEvent generates an event related to a container with only the default attributes.
func (daemon *Daemon) LogContainerEvent(container *container.Container, action string) {
	daemon.LogContainerEventWithAttributes(container, action, map[string]string{})
}

// LogContainerEventWithAttributes generates an event related to a container with specific given attributes.
func (daemon *Daemon) LogContainerEventWithAttributes(container *container.Container, action string, attributes map[string]string) {
	copyAttributes(attributes, container.Config.Labels)
	if container.Config.Image != "" {
		attributes["image"] = container.Config.Image
	}
	attributes["name"] = strings.TrimLeft(container.Name, "/")

	actor := events.Actor{
		ID:         container.ID,
		Attributes: attributes,
	}
	daemon.EventsService.Log(action, events.ContainerEventType, actor)
}

// LogPluginEvent generates an event related to a plugin with only the default attributes.
func (daemon *Daemon) LogPluginEvent(pluginID, refName, action string) {
	daemon.LogPluginEventWithAttributes(pluginID, refName, action, map[string]string{})
}

// LogPluginEventWithAttributes generates an event related to a plugin with specific given attributes.
func (daemon *Daemon) LogPluginEventWithAttributes(pluginID, refName, action string, attributes map[string]string) {
	attributes["name"] = refName
	actor := events.Actor{
		ID:         pluginID,
		Attributes: attributes,
	}
	daemon.EventsService.Log(action, events.PluginEventType, actor)
}

// LogVolumeEvent generates an event related to a volume.
func (daemon *Daemon) LogVolumeEvent(volumeID, action string, attributes map[string]string) {
	actor := events.Actor{
		ID:         volumeID,
		Attributes: attributes,
	}
	daemon.EventsService.Log(action, events.VolumeEventType, actor)
}

// LogNetworkEvent generates an event related to a network with only the default attributes.
func (daemon *Daemon) LogNetworkEvent(nw libnetwork.Network, action string) {
	daemon.LogNetworkEventWithAttributes(nw, action, map[string]string{})
}

// LogNetworkEventWithAttributes generates an event related to a network with specific given attributes.
func (daemon *Daemon) LogNetworkEventWithAttributes(nw libnetwork.Network, action string, attributes map[string]string) {
	attributes["name"] = nw.Name()
	attributes["type"] = nw.Type()
	actor := events.Actor{
		ID:         nw.ID(),
		Attributes: attributes,
	}
	daemon.EventsService.Log(action, events.NetworkEventType, actor)
}

// LogDaemonEventWithAttributes generates an event related to the daemon itself with specific given attributes.
func (daemon *Daemon) LogDaemonEventWithAttributes(action string, attributes map[string]string) {
	if daemon.EventsService != nil {
		if name := hostName(); name != "" {
			attributes["name"] = name
		}
		daemon.EventsService.Log(action, events.DaemonEventType, events.Actor{
			ID:         daemon.id,
			Attributes: attributes,
		})
	}
}

// SubscribeToEvents returns the currently record of events, a channel to stream new events from, and a function to cancel the stream of events.
func (daemon *Daemon) SubscribeToEvents(since, until time.Time, filter filters.Args) ([]events.Message, chan interface{}) {
	ef := daemonevents.NewFilter(filter)
	return daemon.EventsService.SubscribeTopic(since, until, ef)
}

// UnsubscribeFromEvents stops the event subscription for a client by closing the
// channel where the daemon sends events to.
func (daemon *Daemon) UnsubscribeFromEvents(listener chan interface{}) {
	daemon.EventsService.Evict(listener)
}

// copyAttributes guarantees that labels are not mutated by event triggers.
func copyAttributes(attributes, labels map[string]string) {
	if labels == nil {
		return
	}
	for k, v := range labels {
		attributes[k] = v
	}
}

func eventTimestamp(meta swarmapi.Meta, action swarmapi.WatchActionKind) time.Time {
	var eventTime time.Time
	switch action {
	case swarmapi.WatchActionKindCreate:
		eventTime, _ = gogotypes.TimestampFromProto(meta.CreatedAt)
	case swarmapi.WatchActionKindUpdate:
		eventTime, _ = gogotypes.TimestampFromProto(meta.UpdatedAt)
	case swarmapi.WatchActionKindRemove:
		// There is no timestamp from store message for remove operations.
		// Use current time.
		eventTime = time.Now()
	}
	return eventTime
}
