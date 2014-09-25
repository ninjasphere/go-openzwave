package openzwave

type Event interface {
	GetNode() Node
}

type nodeEvent struct {
	node Node
}

type NodeAvailable struct {
	nodeEvent
}

type NodeChanged struct {
	nodeEvent
}

type NodeUnavailable struct {
	nodeEvent
}

func (event nodeEvent) String() string {
	return event.node.(*node).String()
}

func (event nodeEvent) GetNode() Node {
	return event.node.(*node)
}
