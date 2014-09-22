package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ninjasphere/go-openzwave/CODE"
	"github.com/ninjasphere/go-openzwave/NT"
)

// The type of notifications received via the API's Notifications() channel.
type Notification struct {
	cRef *C.Notification
}

// Converts the notification into a string representation.
func (self Notification) String() string {
	return fmt.Sprintf(
		"Notification["+
			"node=0x%08x:0x%02x, "+
			"notificationType=%s/%s, "+
			"valueId=%s, "+
			"value=%s]",
		self.cRef.nodeId.homeId,
		self.cRef.nodeId.nodeId,
		NT.ToEnum(int(self.cRef.notificationType)),
		CODE.ToEnum(int(self.cRef.notificationCode)),
		self.GetValueID(),
		self.GetValue())
}

func (apiNotification *Notification) Free() {
	C.freeNotification(apiNotification.cRef)
}

func (notification *Notification) GetValueID() *ValueID {
	return &ValueID{notification.cRef.valueId}
}

func (notification *Notification) GetValue() *Value {
	return &Value{notification.cRef.value}
}

// given a C notification, return the equivalent Go notification
func asNotification(cRef *C.Notification) *Notification {
	return (*Notification)(unsafe.Pointer(cRef.goRef))
}

//export newGoNotification
func newGoNotification(cRef *C.Notification) unsafe.Pointer {
	goRef := &Notification{cRef}
	cRef.goRef = unsafe.Pointer(goRef)
	return cRef.goRef
}
