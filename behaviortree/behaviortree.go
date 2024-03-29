package behaviortree

import (
	"errors"
	"fmt"
)

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
	// Node represents an node in a tree, that can be ticked
	Node func() (Tick, []Node)

	// Tick represents the logic for a node, which may or may not be stateful
	Tick func(children []Node) (Status, error)

	// Status is a type with three valid values, Running, Success, and Failure, the three possible states for BTs
	Status int
)

// New constructs a new behavior tree and is equivalent to NewNode with vararg support for less indentation
func New(tick Tick, children ...Node) Node {
	return factory(tick, children)
}

func NewNode(tick Tick, children []Node) Node { return factory(tick, children) }

var factory = func(tick Tick, children []Node) (node Node) {
	var (
		frame *Frame
	)
	if v := make([]uintptr, 1); runtimeCallers(3, v[:]) >= 1 {
		if v, _ := runtimeCallersFrames(v).Next(); v.PC != 0 {
			frame = &Frame{
				PC:       v.PC,
				Function: v.Function,
				File:     v.File,
				Line:     v.Line,
				Entry:    v.Entry,
			}
		}
	}
	node = func() (Tick, []Node) {
		if frame != nil {
			node.valueHandle(func(key interface{}) (interface{}, bool) {
				if key != (vkFrame{}) {
					return nil, false
				}
				frame := *frame
				return &frame, true
			})
		}
		return tick, children
	}
	return
}

func (n Node) Tick() (Status, error) {
	if n == nil {
		return Failure, errors.New("behaviortree.Node cannot tick a nil node")
	}
	tick, children := n()
	if tick == nil {
		return Failure, errors.New("behaviortree.Node cannot tick a node with a nil tick")
	}
	return tick(children)
}

func (s Status) Status() Status {
	switch s {
	case Running:
		return Running
	case Success:
		return Success
	default:
		return Failure
	}
}

func (s Status) String() string {
	switch s {
	case Running:
		return `running`
	case Success:
		return `success`
	case Failure:
		return `failure`
	default:
		return fmt.Sprintf("unknown status (%d)", s)
	}
}
