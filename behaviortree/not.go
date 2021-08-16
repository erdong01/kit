package behaviortree

func Not(tick Tick) Tick {
	if tick == nil {
		return nil
	}
	return func(children []Node) (Status, error) {
		status, err := tick(children)
		if err != nil {
			return Failure, err
		}
		switch status {
		case Running:
			return Running, nil
		case Failure:
			return Success, nil
		default:
			return Failure, nil
		}
	}
}
