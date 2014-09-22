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
	"github.com/ninjasphere/go-openzwave/VT"
)

// The type of notifications received via the API's Notifications() channel.
type Notification struct {
	notification *C.Notification
}

// Converts the notification into a string representation.
func (self Notification) String() string {
	return fmt.Sprintf(
		"Notification["+
			"node=0x%08x:0x%02x, "+
			"notificationType=%s/%s, "+
			"valueType=%s, "+
			"valueId=0x%08x]",
		self.notification.nodeId.homeId,
		self.notification.nodeId.nodeId,
		NT.ToEnum(int(self.notification.notificationType)),
		CODE.ToEnum(int(self.notification.notificationCode)),
		VT.ToEnum(int(self.notification.valueId.valueType)),
		self.notification.valueId.id)
}

func (self api) FreeNotification(apiNotification Notification) {
	C.freeNotification(apiNotification.notification)
}

