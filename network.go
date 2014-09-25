package openzwave

import (
	"github.com/ninjasphere/go-openzwave/NT"
)

const (
	MAX_NODES = 232 // max set by ZWave protocol
)

// represents a single Zwave network
type Network interface {
	// the identifier of the home network
	GetHomeId() uint32
}

type network struct {
	homeId uint32
	nodes  map[uint8]*node
}

func newNetwork(homeId uint32) *network {
	return &network{homeId, make(map[uint8]*node)}
}

func (self *network) GetHomeId() uint32 {
	return self.homeId
}

func (self *network) Notify(api API, nt Notification) {
	notificationType := nt.GetNotificationType()
	switch notificationType.Code {

	// network level events
	case NT.DRIVER_READY:
	case NT.DRIVER_RESET:
		// reset network object to reset state
		self.reset()
		break

	// group associations
	case NT.GROUP:
		// not much to do here unless we end up needing to configure group configurations
		// in order to rescue a broken ninja device.

	case NT.AWAKE_NODES_QUERIED:
	case NT.ALL_NODES_QUERIED_SOME_DEAD:
	case NT.ALL_NODES_QUERIED:
		unhandled(api, nt)
		break
		// move network into running state

	// notifications
	case NT.NOTIFICATION:
	default:
		node := nt.GetNode()
		if node.GetId() < MAX_NODES {
			self.handleNodeEvent(api, nt.(*notification), self, self.takeNode(nt.(*notification)))
		} else {
			unexpected(api, nt)
		}
	}
}

func (self *network) handleNodeEvent(api API, nt *notification, net *network, nodeV *node) {

	notificationType := nt.cRef.notificationType
	id := (uint8)(nodeV.cRef.nodeId.nodeId)

	n, ok := net.nodes[id]

	switch notificationType {
	case NT.NODE_REMOVED:
		if ok {
			n.Notify(api, Notification(nt))
			delete(net.nodes, id)
		}
		break

	case NT.NODE_NEW:
	case NT.NODE_ADDED:
		if !ok {
			net.nodes[id] = nodeV
		}

	//
	// node level events
	//
	case NT.NODE_NAMING:
	case NT.NODE_PROTOCOL_INFO:
		// log the related information for diagnostics purposes

	case NT.ESSENTIAL_NODE_QUERIES_COMPLETE:
	case NT.NODE_QUERIES_COMPLETE:
		// move the node into the initialized state
		// begin admission processing for the node

	case NT.VALUE_ADDED:
	case NT.VALUE_REMOVED:
	case NT.VALUE_CHANGED:
	case NT.VALUE_REFRESHED:
		// update node state
		// generate a node changed event

	default:
		// network or node level events
		n.Notify(api, Notification(nt))
		break

	}
}

func (self *network) reset() {
	self.nodes = make(map[uint8]*node)
}

func (self *network) takeNode(nt *notification) *node {
	id := uint8(nt.cRef.node.nodeId.nodeId)
	n, ok := self.nodes[id]
	if !ok {
		n = nt.swapNodeImpl(nil)
		self.nodes[id] = n
	} else {
		nt.swapNodeImpl(n)
	}
	return n
}
