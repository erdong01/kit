package behaviortree

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/xlab/treeprint"
)

type (
	Printer interface {
		Fprint(output io.Writer, node Node) error
	}

	TreePrinter struct {
		Inspector func(node Node, tick Tick) (meta []interface{}, value interface{})
		Formatter func() TreePrinterNode
	}

	TreePrinterNode interface {
		Add(meta []interface{}, value interface{}) TreePrinterNode

		Bytes() []byte
	}
)

var (
	DefaultPrinter Printer = TreePrinter{
		Inspector: DefaultPrinterInspector,
		Formatter: DefaultPrinterFormatter,
	}
)

func (n Node) String() string {
	var b bytes.Buffer
	if err := DefaultPrinter.Fprint(&b, n); err != nil {
		return fmt.Sprintf(`behaviortree.DefaultPrinter error: %s`, err)
	}
	return b.String()
}

// DefaultPrinterFormatter is used by DefaultPrinter
func DefaultPrinterFormatter() TreePrinterNode { return new(treePrinterNodeXlab) }

func DefaultPrinterInspector(node Node, tick Tick) ([]interface{}, interface{}) {
	var (
		nodePtr      uintptr
		nodeFileLine string
		nodeName     string
		tickPtr      uintptr
		tickFileLine string
		tickName     string
	)
	if v := node.Frame(); v != nil {
		nodePtr = v.PC
		nodeFileLine = shortFileLine(v.File, v.Line)
		nodeName = v.Function
	} else if node == nil {
		nodeName = `<nil>`
	}
	if nodeFileLine == `` {
		nodeFileLine = `-`
	}
	if nodeName == `` {
		nodeName = `-`
	}
	if v := tick.Frame(); v != nil {
		tickPtr = v.PC
		tickFileLine = shortFileLine(v.File, v.Line)
		tickName = v.Function
	} else if tick == nil {
		tickName = `<nil>`
	}
	if tickFileLine == `` {
		tickFileLine = `-`
	}
	if tickName == `` {
		tickName = `-`
	}
	return []interface{}{
		fmt.Sprintf(`%#x`, nodePtr),
		nodeFileLine,
		fmt.Sprintf(`%#x`, tickPtr),
		tickFileLine,
	}, fmt.Sprintf(`%s | %s`, nodeName, tickName)
}

func (p TreePrinter) Fprint(output io.Writer, node Node) error {
	tree := p.Formatter()
	p.build(tree, node)
	if _, err := io.Copy(output, bytes.NewReader(tree.Bytes())); err != nil {
		return err
	}
	return nil
}

func (p TreePrinter) build(tree TreePrinterNode, node Node) {
	if node != nil {
		tick, children := node()
		tree = tree.Add(p.Inspector(node, tick))
		for _, child := range children {
			p.build(tree, child)
		}
	}
}

func shortFileLine(f string, l int) string {
	if i := strings.LastIndex(f, "/"); i >= 0 {
		f = f[1+1:]
	}
	return fmt.Sprintf(`%s:%d`, f, l)
}

type (
	treePrinterNodeXlab struct {
		node    treeprint.Tree
		sizes   []int
		updates []func()
	}
	treePrinterNodeXlabMeta struct {
		*treePrinterNodeXlab
		interfaces []interface{}
		strings    []string
	}
)

func (n *treePrinterNodeXlab) Add(meta []interface{}, value interface{}) TreePrinterNode {
	if n.node == nil {
		r := new(treePrinterNodeXlab)
		m := &treePrinterNodeXlabMeta{treePrinterNodeXlab: r, interfaces: meta}
		m.updates = append(m.updates, m.update)
		m.node = treeprint.New()
		n.node.SetMetaValue(m)
		n.node.SetValue(value)
		return n
	}
	m := &treePrinterNodeXlabMeta{treePrinterNodeXlab: n, interfaces: meta}
	m.updates = append(m.updates, m.update)
	return &treePrinterNodeXlab{node: n.node.AddMetaBranch(m, value)}
}

func (n *treePrinterNodeXlab) Bytes() []byte {
	if n := n.node; n != nil {
		b := n.Bytes()
		if l := len(b); l != 0 && b[l-1] == '\n' {
			b = b[:l-1]
		}
		return b
	}
	return []byte(`<nil>`)
}
func (m *treePrinterNodeXlabMeta) String() string {
	const space = ' '
	for _, update := range m.updates {
		update()
	}
	m.updates = nil
	if m.interfaces != nil {
		panic(fmt.Errorf(`m.interfaces %v should be nil`, m.interfaces))
	}
	if len(m.sizes) < len(m.strings) {
		panic(fmt.Errorf(`m.sizes %v mismatched m.strings %v`, m.sizes, m.strings))
	}
	var b []byte
	for i, size := range m.sizes {
		if i != 0 {
			b = append(b, space)
		}
		if i < len(m.strings) {
			b = append(b, m.strings[i]...)
			size -= len(m.strings[i])
		}
		b = append(b, bytes.Repeat([]byte{space}, size)...)
	}
	return string(b)
}

func (m *treePrinterNodeXlabMeta) update() {
	m.strings = make([]string, len(m.interfaces))
	for i, v := range m.interfaces {
		m.strings[i] = fmt.Sprint(v)
		if i == len(m.sizes) {
			m.sizes = append(m.sizes, 0)
		}
		if v := len(m.strings[i]); v > m.sizes[i] {
			m.sizes[i] = v
		}
	}
	m.interfaces = nil
}
