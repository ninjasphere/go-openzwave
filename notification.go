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
func (self *notification) String() string {
	return fmt.Sprintf(
		"Notification["+
			"notificationType=%v/%v, "+
			"node=%v, "+
			"value=%v]",
		NT.ToEnum(int(self.cRef.notificationType)),
		CODE.ToEnum(int(self.cRef.notificationCode)),
		self.GetNode(),
		self.GetValue())
}

func (self *notification) free() {
	C.freeNotification(self.cRef)
	if self.node != nil {
		self.node.free()
	}
	if self.value != nil {
		self.value.free()
	}
}

func (self *notification) GetValue() Value {
	return self.value
}

func (self *notification) GetNode() Node {
	return self.node
}

func (self *notification) GetNotificationType() *NT.Enum {
	return NT.ToEnum(int(self.cRef.notificationType))
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
func (self *notification) swapNodeImpl(existing *node) *node {
	if existing != nil {
		// then swap the cRef pointers
		swap := self.node.cRef
		self.node.cRef = existing.cRef
		existing.cRef = swap
	} else {
		existing = self.node
		self.node = nil
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
func (self *notification) swapValueImpl(existing *value) *value {
	if existing != nil {
		// then swap the cRef pointers
		swap := self.value.cRef
		self.value.cRef = existing.cRef
		existing.cRef = swap
	} else {
		existing = self.value
		self.value = nil
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
