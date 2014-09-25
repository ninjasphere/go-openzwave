package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ninjasphere/go-openzwave/NT"
)

type Node interface {
	GetHomeId() uint32
	GetId() uint8
}

type node struct {
	cRef    *C.Node
	classes map[uint8]*valueClass
}

type valueClass struct {
	commandClass uint8
	instances    map[uint8]*valueInstance
}

type valueInstance struct {
	instance uint8
	values   map[uint8]*value
}

func (self *node) String() string {
	cRef := self.cRef

	return fmt.Sprintf(
		"Node["+
			"homeId=0x%08x, "+
			"nodeId=%03d, "+
			"basicType=%02x, "+
			"genericType=%02x, "+
			"specificType=%02x, "+
			"nodeType='%s', "+
			"manufacturerName='%s', "+
			"productName='%s', "+
			"location='%s', "+
			"manufacturerId=%s, "+
			"productType=%s, "+
			"productId=%s]",
		uint32(cRef.nodeId.homeId),
		uint8(cRef.nodeId.nodeId),
		uint8(cRef.basicType),
		uint8(cRef.genericType),
		uint8(cRef.specificType),
		C.GoString(cRef.nodeType),
		C.GoString(cRef.manufacturerName),
		C.GoString(cRef.productName),
		C.GoString(cRef.location),
		C.GoString(cRef.manufacturerId),
		C.GoString(cRef.productType),
		C.GoString(cRef.productId))
}

// convert a reference from the C Node to the Go Node
func asNode(cRef *C.Node) Node {
	return Node((*node)(cRef.goRef))
}

//export newGoNode
func newGoNode(cRef *C.Node) unsafe.Pointer {
	goRef := &node{cRef, make(map[uint8]*valueClass)}
	cRef.goRef = unsafe.Pointer(goRef)
	return cRef.goRef
}

func (self *node) GetHomeId() uint32 {
	return uint32(self.cRef.nodeId.homeId)
}

func (self *node) GetId() uint8 {
	return uint8(self.cRef.nodeId.nodeId)
}

func (self *node) Notify(api API, nt Notification) {
	notificationType := nt.GetNotificationType()
	switch notificationType.Code {
	case NT.VALUE_REMOVED:
		self.removeValue(nt.(*notification))
		break

	case NT.VALUE_ADDED:
	case NT.VALUE_CHANGED:
	case NT.VALUE_REFRESHED:
		self.takeValue(nt.(*notification))
		break

	case NT.NODE_NAMING:
	case NT.NODE_PROTOCOL_INFO:
		// log the related information for diagnostics purposes

	case NT.ESSENTIAL_NODE_QUERIES_COMPLETE:
	case NT.NODE_QUERIES_COMPLETE:
		// move the node into the initialized state
		// begin admission processing for the node
		break
	}
}

// take the value structure from the notification
func (self *node) takeValue(nt *notification) *value {
	commandClassId := (uint8)(nt.cRef.value.valueId.commandClassId)
	instanceId := (uint8)(nt.cRef.value.valueId.instance)
	index := (uint8)(nt.cRef.value.valueId.index)

	instance := self.createOrGetInstance(commandClassId, instanceId)
	v, ok := instance.values[index]
	if !ok {
		v = nt.swapValueImpl(nil)
		instance.values[index] = v
	} else {
		nt.swapValueImpl(v)
	}
	return v
}

func (self *node) createOrGetInstance(commandClassId uint8, instanceId uint8) *valueInstance {
	class, ok := self.classes[commandClassId]
	if !ok {
		class = &valueClass{commandClassId, make(map[uint8]*valueInstance)}
		self.classes[commandClassId] = class
	}
	instance, ok := class.instances[instanceId]
	if !ok {
		instance = &valueInstance{instanceId, make(map[uint8]*value)}
		class.instances[instanceId] = instance
	}
	return instance
}

func (self *node) removeValue(nt *notification) {
	commandClassId := (uint8)(nt.cRef.value.valueId.commandClassId)
	instanceId := (uint8)(nt.cRef.value.valueId.instance)
	index := (uint8)(nt.cRef.value.valueId.index)

	class, ok := self.classes[commandClassId]
	if !ok {
		return
	}

	instance, ok := class.instances[instanceId]
	if !ok {
		return
	}

	value, ok := instance.values[index]
	_ = value

	if !ok {
		return
	} else {
		delete(instance.values, index)
		if len(instance.values) == 0 {
			delete(class.instances, instanceId)
			if len(class.instances) == 0 {
				delete(self.classes, commandClassId)
			}
		}
	}

}
