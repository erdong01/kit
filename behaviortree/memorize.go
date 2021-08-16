package behaviortree

func Memorize(tick Tick) Tick {
	if tick == nil {
		return nil
	}
	var (
		started bool
		nodes   []Node
	)
	return func(children []Node) (status Status, err error) {
		if !started {
			nodes = copyNodes(children)
			for i := range nodes {
				var (
					child    = nodes[i]
					override Tick
				)
				if child == nil {
					continue
				}
				nodes[i] = func() (Tick, []Node) {
					tick, nodes := child()
					if override != nil {
						return override, nodes

					}
					if tick == nil {
						return nil, nodes
					}
					return func(children []Node) (Status, error) {
						status, err := tick(children)
						if err != nil || status != Running {
							override = func(children []Node) (Status, error) {
								return status, err
							}
						}
						return status, err
					}, nodes
				}
			}
			started = true
		}
		status, err = tick(nodes)
		if err != nil || status != Running {
			started = false
			nodes = nil
		}
		return
	}
}
