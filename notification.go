package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"

	"github.com/ninjasphere/go-openzwave/CODE"
	"github.com/ninjasphere/go-openzwave/NT"
)

// The type of notifications received from the API
type Notification interface {
	GetNode() Node
	GetNotificationType() *NT.Enum
}

// The type of notifications received via the API's Notifications() channel.
type notification struct {
	cRef  *C.Notification
	node  *node  // should be freed by the receiver, iff it is not null
	value *value // should be freed by the receiver, iff it is not null
}

// Converts the notification into a string representation.
func (n *notification) String() string {
	return fmt.Sprintf(
		"Notification["+
			"notificationType=%v/%v, "+
			"node=%v, "+
			"value=%v]",
		NT.ToEnum(int(n.cRef.notificationType)),
		CODE.ToEnum(int(n.cRef.notificationCode)),
		n.GetNode(),
		n.GetValue())
}

func (n *notification) free() {
	C.freeNotification(n.cRef)
	if n.node != nil {
		n.node.free()
	}
	if n.value != nil {
		n.value.free()
	}
}

func (n *notification) GetValue() Value {
	return n.value
}

func (n *notification) GetNode() Node {
	return n.node
}

func (n *notification) GetNotificationType() *NT.Enum {
	return NT.ToEnum(int(n.cRef.notificationType))
}

func newGoNotification(cRef *C.Notification) *notification {
	result := &notification{cRef, newGoNode(cRef.node), newGoValue(cRef.value)}

	// transfer ownership of C structure to the go object
	result.cRef.value = nil
	result.cRef.node = nil
	return result
}

//
// Swap the cRef of the receiver's node with cRef of the specified node.
//
// The intent is to swap the *C.Node reachable from 'existing' with the *C.Node
// reachable from the receiver so that existing gets the fresh object.
//
// If there is no existing object then we steal the whole go object from the
// receiver.
//
func (n *notification) swapNodeImpl(existing *node) *node {
	if existing != nil {
		// then swap the cRef pointers
		swap := n.node.cRef
		n.node.cRef = existing.cRef
		existing.cRef = swap
	} else {
		existing = n.node
		n.node = nil
	}
	return existing
}

//
// Swap the cRef of the receiver's node with cRef of the specified node.
//
// The intent is to update an existing go representation of a node with
// the latest *C.Value from the notification and then recycle
// the old *C.Value by attaching it to the notification where it will then
// be freed.
//
func (n *notification) swapValueImpl(existing *value) *value {
	if existing != nil {
		// then swap the cRef pointers
		swap := n.value.cRef
		n.value.cRef = existing.cRef
		existing.cRef = swap
	} else {
		existing = n.value
		n.value = nil
	}
	return existing
}

//
// called for unexpected notifications.
//
func unexpected(api API, notification Notification) {
	api.Logger().Warningf("unexpected notification received %v]\n", notification)
}

//
// called for expected notifications that are not handled
//
func unhandled(api API, notification Notification) {
	// api.Logger().Debugf("unhandled notification received %v]\n", notification)
}
