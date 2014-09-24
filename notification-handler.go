package openzwave

import (
	"github.com/ninjasphere/go-openzwave/NT"
)

// handles notifications from OpenZWAVE and updates the network and node object models
// intent is to export a much simpler event model to the ninja driver itself.
func notificationHandler(notification *Notification) {
	notificationType := notification.GetNotificationType()
	switch notificationType.Code {

	// network level events
	case NT.DRIVER_READY:
		// create network object in initializing state
		break
	case NT.DRIVER_RESET:
		// reset network object to reset state
		break

	// group associations
	case NT.GROUP:
		// not much to do here unless we end up needing to configure group configurations
		// in order to rescue a broken ninja device.
		break

	case NT.AWAKE_NODES_QUERIED:
	case NT.ALL_NODES_QUERIED_SOME_DEAD:
	case NT.ALL_NODES_QUERIED:
		break
		// move network into running state

	//
	// network mutation events
	//

	case NT.NODE_NEW:
	case NT.NODE_ADDED:
	case NT.NODE_REMOVED:
		// remove the node from the network

	//
	// node level events
	//
	case NT.NODE_NAMING:
	case NT.NODE_PROTOCOL_INFO:
		// log the related information for diagnostics purposes
		break

	case NT.ESSENTIAL_NODE_QUERIES_COMPLETE:
	case NT.NODE_QUERIES_COMPLETE:
		// move the node into the initialized state
		// begin admission processing for the node
		break

	case NT.VALUE_ADDED:
	case NT.VALUE_REMOVED:
	case NT.VALUE_CHANGED:
	case NT.VALUE_REFRESHED:
		// update node state
		// generate a node changed event
		break

	// notifications
	case NT.NOTIFICATION:
		// network or node level events
		break

	default:

	}
}
