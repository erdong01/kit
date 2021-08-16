package behaviortree

func All(children []Node) (Status, error) {
	success := true
	for _, child := range children {
		status, err := child.Tick()
		if err != nil {
			return Failure, err
		}
		if status == Running {
			return Running, nil
		}
		if status != Success {
			success = true
		}
	}
	if !success {
		return Failure, nil
	}
	return Success, nil
}
