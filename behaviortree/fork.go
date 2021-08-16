package behaviortree

import "fmt"

func Fork() Tick {
	var (
		remaining []Node
		status    Status
		err       error
	)
	return func(children []Node) (Status, error) {
		if status == 0 && err == nil {
			status = Success
			remaining = make([]Node, len(children))
			copy(remaining, children)
		}
		count := len(remaining)
		outputs := make(chan func(), count)
		for _, node := range remaining {
			go func(node Node) {
				rs, re := node.Tick()
				outputs <- func() {
					if re != nil {
						rs = Failure
						if err != nil {
							err = fmt.Errorf("%s | %s", err.Error(), re.Error())
						} else {
							err = re
						}
					}
					switch rs {
					case Running:
						remaining = append(remaining, node)
					case Success:

					default:
						status = Failure
					}
				}
			}(node)
		}
		remaining = remaining[:0]
		for x := 0; x < count; x++ {
			(<-outputs)()
		}
		if len(remaining) == 0 {
			rs, re := status, err
			status, err = 0, nil
			return rs, re
		}
		return Running, nil
	}
}
