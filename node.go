package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"

	"github.com/ninjasphere/go-openzwave/NT"
)

type state int

const (
	STATE_INIT  state = iota
	STATE_READY       = iota
)

type Node interface {
	GetHomeId() uint32
	GetId() uint8

	GetDevice() Device

	GetProductId() *ProductId
	GetProductDescription() *ProductDescription
	GetNodeName() string

	GetValue(commandClassId uint8, instanceId uint8, index uint8) Value
}

type ProductId struct {
	ManufacturerId string
	ProductId      string
}

type ProductDescription struct {
	ManufacturerName string
	ProductName      string
	ProductType      string
}

type node struct {
	cRef    *C.Node
	classes map[uint8]*valueClass
	state   state
	device  Device
}

type valueClass struct {
	commandClass uint8
	instances    map[uint8]*valueInstance
}

type valueInstance struct {
	instance uint8
	values   map[uint8]*value
}

func newGoNode(cRef *C.Node) *node {
	return &node{cRef, make(map[uint8]*valueClass), STATE_INIT, nil}
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

func (self *node) GetHomeId() uint32 {
	return uint32(self.cRef.nodeId.homeId)
}

func (self *node) GetId() uint8 {
	return uint8(self.cRef.nodeId.nodeId)
}

func (self *node) notify(api *api, nt *notification) {

	var event Event

	notificationType := nt.cRef.notificationType
	switch notificationType {
	case NT.NODE_REMOVED:
		event = &NodeUnavailable{nodeEvent{self}}
		if self.device != nil {
			self.device.NodeRemoved()
		}
		api.notifyEvent(event)
		// TODO: free the C structure.
		break

	case NT.VALUE_REMOVED:
		self.removeValue(nt)
		break

	case NT.ESSENTIAL_NODE_QUERIES_COMPLETE,
		NT.NODE_QUERIES_COMPLETE:
		// move the node into the initialized state
		// begin admission processing for the node

		switch self.state {
		case STATE_INIT:
			self.state = STATE_READY

			event = &NodeAvailable{nodeEvent{self}}
			//
			// Use a callback to construct the device for this node, then
			// pass the event to the device.
			//

			self.device = api.deviceFactory(api, self)
			self.device.NodeAdded()

			break
		default:
			event = &NodeChanged{nodeEvent{self}}
			self.device.NodeChanged()
			//
			// Pass the event to the node.
			//
		}
		api.notifyEvent(event)
		break

	case NT.VALUE_ADDED,
		NT.VALUE_CHANGED,
		NT.VALUE_REFRESHED:
		self.takeValue(nt)
		break

	case NT.NODE_NAMING,
		NT.NODE_PROTOCOL_INFO:
		// log the related information for diagnostics purposes

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

func (self *node) GetValue(commandClassId uint8, instanceId uint8, index uint8) Value {
	var v *value
	class, ok := self.classes[commandClassId]
	if ok {
		instance, ok := class.instances[instanceId]
		if ok {
			v, ok = instance.values[index]
		}
	}
	if ok {
		return v
	} else {
		return &missingValue{} // accessor that does nothing
	}
	return v
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
		// TODO: free the C structure
		delete(instance.values, index)
		if len(instance.values) == 0 {
			delete(class.instances, instanceId)
			if len(class.instances) == 0 {
				delete(self.classes, commandClassId)
			}
		}
	}

}

func (self *node) GetDevice() Device {
	return self.device
}

func (self *node) GetProductId() *ProductId {
	return &ProductId{C.GoString(self.cRef.manufacturerId), C.GoString(self.cRef.productId)}
}

func (self *node) GetProductDescription() *ProductDescription {
	return &ProductDescription{
		C.GoString(self.cRef.manufacturerName),
		C.GoString(self.cRef.productName),
		C.GoString(self.cRef.productType)}
}

func (self *node) GetNodeName() string {
	return C.GoString(self.cRef.nodeName)
}
