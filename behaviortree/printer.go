package behaviortree

import "io"

type (
	Printer interface {
		Fprint(output io.Writer, node Node) error
	}

	TreePrinter struct {
		Inspector func(node Node, tick Tick) (meta []interface{}, value interface{})
		Formatter func() TreePrinterNode
	}
)
