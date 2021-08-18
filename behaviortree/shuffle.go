package behaviortree

import "math/rand"

func Shuffle(tick Tick, source rand.Source) Tick {
	if tick == nil {
		return nil
	}

	if source == nil {
		source = defaultSource{}
	}
	return func(children []Node) (Status, error) {
		children = copyNodes(children)
		rand.New(source).Shuffle(len(children), func(i, j int) {
			children[i], children[j] = children[j], children[i]

		})
		return tick(children)
	}
}

type defaultSource struct{ rand.Source }

func (d defaultSource) Int63() int64 {
	return rand.Int63()
}
