package behaviortree

const (
	_Status = iota
	// Running indicates that the Tick for a given Node is currently running
	Running
	// Success indicates that the Tick for a given Node completed successfully
	Success
	// Failure indicates that the Tick for a given Node failed to complete successfully
	Failure
)

type (
	Node func(Tick, []Node)

	Tick func(children []Node) (Status, error)

	Status int
)
