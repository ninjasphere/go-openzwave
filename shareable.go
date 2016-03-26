package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
// #include "api.h"
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// This module fixes an issue revealed when Go 1.6 tightened up the rules
// about sharing of Go pointers with C code. https://github.com/golang/go/issues/12416
//
// Now we register a reference with a map and get a pointer to a structure in the C heap on return.
// We use the integer recorded in this structure to map back into the Go world.
//
// That way there are O(1) marshaling operations, and Go pointer is never shared with or
// dereferenced by the C world
//
var shared = map[int]Shareable{}
var sharedCount = 0
var mu sync.RWMutex

type Shareable interface {
	C() unsafe.Pointer
	Go() interface{}
}

type shareable struct {
	cref     *C.shareable
	goObject interface{}
}

func (s *shareable) init(goObject interface{}) {
	mu.Lock()
	defer mu.Unlock()

	sharedCount++
	s.cref = C.newShareable(C.int(sharedCount))
	shared[int(s.cref.sharedIndex)] = s
	s.goObject = goObject
}

func (s *shareable) destroy() {
	if s.cref != nil {
		mu.Lock()
		defer mu.Unlock()
		delete(shared, int(s.cref.sharedIndex))
		C.free(s.C())
		s.cref = nil
	}
}

func (s *shareable) C() unsafe.Pointer {
	return unsafe.Pointer(s.cref)
}

func (s *shareable) Go() interface{} {
	if s == nil {
		return nil
	} else {
		return s.goObject
	}
}

func unmarshal(c unsafe.Pointer) Shareable {
	mu.RLock()
	defer mu.RUnlock()

	if c == nil {
		return nil
	}

	i := (*C.shareable)(c).sharedIndex
	if s, ok := shared[int(i)]; !ok {
		panic(fmt.Errorf("failure to unmarshal index %d", int(i)))
	} else {
		return s
	}
}
