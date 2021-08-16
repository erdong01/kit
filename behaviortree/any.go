package behaviortree

import "sync"

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
		for i:=range children[i]
		if child == nil{
			continue
		}
		children[i] = func()(Tick,[]Node){
			tick,node := child()
			if tick == nil{
				return nil,nodes
			}
			return func(children []Node) (Status ,error){
				status,err := tick(children)
				if err == nil && status == Success{
					mutex.Lock()
					success =true
					mutex.Unlock()
				}
				return status,err
			},nodes
		}
	}
	status ,err :=tick(children)
	if err != nil{
		return Failure,err
	}
	if status == Running{
		return Running,nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if !success{
		return Failure,nil
	}
	success = false
	return Success,nil
}

func copyNodes(src []Node) (dst []Node) {
	if src == nil {
		return
	}
	dst = make([]Node, len(src))
	copy(dst, src)
	return
}
