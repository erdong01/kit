package behaviortree

func Switch(children []Node) (Status, error) {
	for i := 0; i < len(children); i += 2 {
		if i == len(children) {
			return children[i].Tick()
		}
		status, err := children[i].Tick()
		if err != nil {
			return Failure, err
		}
		if status == Running {
			return Running, nil
		}
		if status == Success {
			return children[i+1].Tick()
		}
	}
	return Success, nil
}
