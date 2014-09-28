package openzwave

// A device is an abstraction used by higher-level parts of the system for the underlying ZWave node.
type Device interface {
	NodeAdded()
	NodeChanged()
	NodeRemoved()
	ValueChanged(Value)
}

type DeviceFactory func(API, Node) Device

type emptyDevice struct {
}

func defaultDeviceFactory(api API, node Node) Device {
	return &emptyDevice{}
}

func (*emptyDevice) NodeAdded() {
}

func (*emptyDevice) NodeChanged() {
}

func (*emptyDevice) NodeRemoved() {
}

func (*emptyDevice) ValueChanged(value Value) {
}
