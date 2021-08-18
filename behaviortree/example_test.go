package behaviortree

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func newExampleCounter() Node {
	var (
		// counter is the shared state used by this example
		counter = 0
		// printCounter returns a node that will print the counter prefixed with the given name then succeed
		printCounter = func(name string) Node {
			return New(
				func(children []Node) (Status, error) {
					fmt.Printf("%s: %d\n", name, counter)
					return Success, nil
				},
			)
		}
		// incrementCounter is a node that will increment counter then succeed
		incrementCounter = New(
			func(children []Node) (Status, error) {
				counter++
				return Success, nil
			},
		)
	)
	return New(
		Selector, // runs each child sequentially until one succeeds (success) or all fail (failure)
		New(
			Sequence, // runs each child in order until one fails (failure) or they all succeed (success)
			New(
				func(children []Node) (Status, error) { // succeeds while counter is less than 10
					if counter < 10 {
						return Success, nil
					}
					return Failure, nil
				},
			),
			incrementCounter,
			printCounter("< 10"),
		),
		New(
			Sequence,
			New(
				func(children []Node) (Status, error) { // succeeds while counter is less than 20
					if counter < 20 {
						return Success, nil
					}
					return Failure, nil
				},
			),
			incrementCounter,
			printCounter("< 20"),
		),
	) // if both children failed (counter is >= 20) the root node will also fail
}

// ExampleNewTickerStopOnFailure_counter demonstrates the use of NewTickerStopOnFailure to implement more complex "run
// to completion" behavior using the simple modular building blocks provided by this package
func TestExampleNewTickerStopOnFailure_counter(t *testing.T) {
	// ticker is what actually runs this example and will tick the behavior tree defined by a given node at a given
	// rate and will stop after the first failed tick or error or context cancel
	ticker := NewTickerStopOnFailure(
		context.Background(),
		time.Millisecond,
		newExampleCounter(),
	)
	// waits until ticker stops, which will be on the first failure of it's root node
	<-ticker.Done()
	// every Tick may return an error which would automatically cause a failure and propagation of the error
	if err := ticker.Err(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)
}
