package behaviortree

func Sequence(children []Node) (Status, error) {

	for _, c := range children {
		status, err := c.Tick()
		if err != nil {
			return Failure, err
		}
		if status == Running {
			return Running, nil
		}
		if status != Success {
			return Failure, nil
		}
	}
	return Success, nil
}
