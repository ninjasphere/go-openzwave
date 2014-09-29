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
	node  *node
	value *value
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

func (apiNotification *notification) free() {
	C.freeNotification(apiNotification.cRef)
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
	return &notification{cRef, newGoNode(cRef.node), newGoValue(cRef.value)}
}

//
// Swap the cRef of the receiver's node with cRef of the specified node.
//
// The intent is to update an existing go representation of a node with
// the latest *C.Node from the notification and then recycle
// the old *C.Node by attaching it to the notification where it will then
// be freed.
//
func (self *notification) swapNodeImpl(existing *node) *node {
	if existing != nil {
		existingGoNode := existing

		// then swap the cRef pointers
		swap := self.cRef.node
		self.cRef.node = existingGoNode.cRef
		existingGoNode.cRef = swap
	} else {
		existing = self.node
		self.cRef.node = nil
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
		existingGoValue := existing

		// then swap the cRef pointers
		swap := self.cRef.value
		self.cRef.value = existingGoValue.cRef
		existingGoValue.cRef = swap
	} else {
		existing = self.value
		self.cRef.value = nil
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
	api.Logger().Debugf("unhandled notification received %v]\n", notification)
}
