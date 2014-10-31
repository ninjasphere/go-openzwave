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

func (nw *network) GetHomeId() uint32 {
	return nw.homeId
}

func (nw *network) notify(api *api, nt *notification) {
	notificationType := nt.GetNotificationType()
	switch notificationType.Code {

	// network level events
	case NT.DRIVER_READY,
		NT.DRIVER_RESET:
		// reset network object to reset state
		nw.reset()
		break

	// group associations
	case NT.GROUP,
		NT.AWAKE_NODES_QUERIED,
		NT.ALL_NODES_QUERIED_SOME_DEAD,
		NT.ALL_NODES_QUERIED:
		unhandled(api, nt)
		break
		// move network into running state

	// notifications
	case NT.NOTIFICATION:
		fallthrough
	default:
		node := nt.GetNode()
		if node.GetId() <= MAX_NODES {
			nw.handleNodeEvent(api, nt, nw.takeNode(nt))
		} else {
			unhandled(api, nt)
		}
	}
}

func (nw *network) handleNodeEvent(api *api, nt *notification, nodeV *node) {

	notificationType := nt.cRef.notificationType
	id := (uint8)(nodeV.cRef.nodeId.nodeId)

	n, ok := nw.nodes[id]

	switch notificationType {
	case NT.NODE_REMOVED:
		if ok {
			delete(nw.nodes, id)
			n.notify(api, nt)
		}
		break

	case NT.ESSENTIAL_NODE_QUERIES_COMPLETE:
	case NT.NODE_QUERIES_COMPLETE:
		// move the node into the initialized state
		// begin admission processing for the node
		// network or node level events
		n.notify(api, nt)
		break

	case NT.NODE_NEW,
		NT.NODE_ADDED:
		if !ok {
			nw.nodes[id] = nodeV
		}
		fallthrough

	//
	// node level events
	//
	case NT.NODE_NAMING,
		NT.NODE_PROTOCOL_INFO,
		NT.VALUE_ADDED,
		NT.VALUE_REMOVED,
		NT.VALUE_CHANGED,
		NT.VALUE_REFRESHED:
		fallthrough
	default:
		// network or node level events
		n.notify(api, nt)
		break

	}
}

func (nw *network) reset() {
	nw.nodes = make(map[uint8]*node)
}

func (nw *network) takeNode(nt *notification) *node {
	id := uint8(nt.node.cRef.nodeId.nodeId)
	n, ok := nw.nodes[id]
	if !ok {
		n = nt.swapNodeImpl(nil)
		nw.nodes[id] = n
	} else {
		nt.swapNodeImpl(n)
	}
	return n
}
