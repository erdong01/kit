package behaviortree

import "sync"

// Any wraps a tick such that non-error non-running statuses will be overridden with a success if at least one child
// succeeded - which is achieved by encapsulation of children, before passing them into the wrapped tick. Nil will be
// returned if tick is nil, and nil children will be passed through as such.
func Any(tick Tick) Tick {
	if tick == nil {
		return nil
	}
	var (
		mutex   sync.Mutex
		success bool
	)
	return func(children []Node) (Status, error) {
		children = copyNodes(children)
		for i := range children {
			child := children[i]
			if child == nil {
				continue
			}
			children[i] = func() (Tick, []Node) {
				tick, nodes := child()
				if tick == nil {
					return nil, nodes
				}
				return func(children []Node) (Status, error) {
					status, err := tick(children)
					if err == nil && status == Success {
						mutex.Lock()
						success = true
						mutex.Unlock()
					}
					return status, err
				}, nodes
			}
		}
		status, err := tick(children)
		if err != nil {
			return Failure, err
		}
		if status == Running {
			return Running, nil
		}
		mutex.Lock()
		defer mutex.Unlock()
		if !success {
			return Failure, nil
		}
		success = false
		return Success, nil
	}
}

func copyNodes(src []Node) (dst []Node) {
	if src == nil {
		return
	}
	dst = make([]Node, len(src))
	copy(dst, src)
	return
}
