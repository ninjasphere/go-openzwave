package openzwave

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
