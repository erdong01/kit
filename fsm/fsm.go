package fsm

type FSM struct {
	Events map[string]IEvent
}

type IEvent interface {
	Run()
}

func NewFSM(initial string, events []IEvent) *FSM {
	f := &FSM{}

	return f
}
