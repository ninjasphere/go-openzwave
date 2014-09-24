package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type Node interface {
	Notifiable
	GetHomeId() uint32
	GetId() uint8
}

type node struct {
	cRef *C.Node
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
	goRef := &node{cRef}
	cRef.goRef = unsafe.Pointer(goRef)
	return cRef.goRef
}

func (self *node) GetHomeId() uint32 {
	return uint32(self.cRef.nodeId.homeId)
}

func (self *node) GetId() uint8 {
	return uint8(self.cRef.nodeId.nodeId)
}

func (self *node) Notify(api API, notification Notification) {
	// TODO
}
