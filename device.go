package openzwave

// A device is an abstraction used by higher-level parts of the system for the underlying ZWave node.
type Device interface {
	// to receive notifications of API events relating to the underlying node
	Notify(API, Event)
}

type DeviceFactory func(API, Node) Device

type emptyDevice struct {
}

func defaultDeviceFactory(api API, node Node) Device {
	return &emptyDevice{}
}

func (*emptyDevice) Notify(api API, event Event) {
}
