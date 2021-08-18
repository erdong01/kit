package behaviortree

import "sync"

type syc struct {
	nodes    []Node
	statuses []Status
	mutex    *sync.Mutex
}

func Sync(nodes []Node) []Node {
	if nodes == nil {
		return nil
	}
	s := &syc{
		nodes:    nodes,
		mutex:    new(sync.Mutex),
		statuses: make([]Status, len(nodes)),
	}
	result := make([]Node, 0, len(nodes))
	for i := range nodes {
		result = append(result, s.node(i))
	}
	return result
}

func (s *syc) running() bool {
	for _, status := range s.statuses {
		if status == Running {
			return true
		}
	}
	return false
}

func (s *syc) node(i int) Node {
	if s.nodes[i] == nil {
		return nil
	}
	return func() (Tick, []Node) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		tick, children := s.nodes[i]()
		if tick == nil {
			return nil, children
		}
		status := s.statuses[i]
		if status != Running && s.running() {
			return func(children []Node) (Status, error) {
				return status, nil
			}, children
		}
		return func(children []Node) (Status, error) {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			status, err := tick(children)
			s.statuses[i] = status
			return status, err
		}, children
	}
}
