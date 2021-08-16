package behaviortree

func Background(tick func() Tick) Tick {
	if tick == nil {
		return nil
	}
	var nodes []Node
	return func(children []Node) (Status, error) {
		for i, node := range nodes {
			status, err := node.Tick()
			if err == nil && status == Running {
				continue
			}
			copy(nodes[i:], nodes[i+1:])
			nodes[len(nodes)-1] = nil
			nodes = nodes[:len(nodes)-1]
			return status, err
		}
		node := NewNode(tick(), children)
		status, err := node.Tick()
		if err != nil || status != Running {
			return status, err
		}
		nodes = append(nodes, node)
		return Running, nil
	}
}
