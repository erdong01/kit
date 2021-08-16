package behaviortree

func Async(tick Tick) Tick {
	if tick == nil {
		return nil
	}
	var done chan struct {
		Status Status
		Error  error
	}
	return func(children []Node) (Status, error) {
		if done == nil {
			done = make(chan struct {
				Status Status
				Error  error
			}, 1)
			go func() {
				var status struct {
					Status Status
					Error  error
				}
				defer func() {
					done <- status
				}()
				status.Status, status.Error = tick(children)
			}()
			return Running, nil
		}

		select {
		case status := <-done:
			done = nil
			return status.Status, status.Error
		default:
			return Running, nil
		}
	}
}
