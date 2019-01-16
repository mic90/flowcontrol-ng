package graph

import (
	"context"
	"fmt"
	"github.com/mic90/flowcontrol-ng/port"
	"github.com/rs/xid"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

type ConcurrentRunner struct {
	state        State
	stop         atomic.Value
	waitGroup    sync.WaitGroup
	error        chan bool
	errorHandler func()
}

func NewConcurrentRunner(errorHandler func()) *ConcurrentRunner {
	return &ConcurrentRunner{StateStopped, atomic.Value{}, sync.WaitGroup{}, make(chan bool), errorHandler}
}

func (runner *ConcurrentRunner) Run(ctx context.Context, nodes map[xid.ID]Node) {
	runner.state = StateStarting
	runner.stop.Store(false)

	// run watcher to check for graph error
	go runner.errorWatcher(ctx)

	for key := range nodes {
		node := nodes[key]
		runner.setupNode(node)
	}

	for key := range nodes {
		node := nodes[key]
		go runner.runNode(ctx, node)
	}
	runner.state = StateStarted
}

func (runner *ConcurrentRunner) Stop() {
	runner.stop.Store(true)
	runner.waitGroup.Wait()
	runner.state = StateStopped
}

func (runner *ConcurrentRunner) Wait() {

}

func (runner *ConcurrentRunner) State() State {
	return runner.state
}

func (runner *ConcurrentRunner) setError(err error) {
	log.Println("Graph error occurred:", err)
	runner.error <- true
	runner.state = StateError
}

func (runner *ConcurrentRunner) errorWatcher(ctx context.Context) {
	select {
	case <-runner.error:
		// set graph internal state to stop
		runner.stop.Store(true)
		// wait for all nodes to stop, then handle error as user requested
		runner.waitGroup.Wait()
		runner.errorHandler()
	case <-ctx.Done():
		runner.Stop()
		log.Println("graph was stopped by context")
	}
}

func (runner *ConcurrentRunner) setupNode(node Node) {
	setupError := node.Setup()
	if setupError != nil {
		log.Println("Couldn't setup node", node.Name())
		runner.setError(setupError)
		return
	}
}

func (runner *ConcurrentRunner) runNode(ctx context.Context, node Node) {
	runner.waitGroup.Add(1)
	defer runner.waitGroup.Done()
	defer runner.stopNode(node)

	processTimer := NewProcessTimer(ctx)
	nodeConnectors := runner.getNodeConnectors(node)
	log.Println("Started node:", node.Name(), node.ID())
	for {
		// stop processing loop
		if runner.stop.Load().(bool) == true {
			break
		}

		// run node processing function, if error occurred stop whole graph
		err := runner.processNode(ctx, processTimer, node, nodeConnectors)
		if err != nil {
			runner.setError(err)
			break
		}
	}
}

func (runner *ConcurrentRunner) processNode(ctx context.Context, timer *ProcessTimer, node Node, connectors []port.Connector) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("unexpected error occured. Reason: %v", recovered)
		}
	}()

	// wait for input values to appear asynchronously
	// to prevent dead-lock in writers
	connectorsCount := len(connectors)
	var wg sync.WaitGroup
	wg.Add(connectorsCount)
	for i := 0; i < connectorsCount; i++ {
		conn := connectors[i]
		go func() {
			defer wg.Done()

			// wait for data based on QOS level for the connection
			conn.Wait()

			// pass data from src to dst port
			err := conn.Trigger()
			if err != nil {
				runner.setError(err)
			}
		}()
	}
	wg.Wait()

	// stop processing loop, double check in case output->input write went wrong
	if runner.stop.Load().(bool) == true {
		return nil
	}

	return node.Process(ctx, timer)
}

func (runner *ConcurrentRunner) getNodeConnectors(node Node) []port.Connector {
	nodeConnectors := make([]port.Connector, 0)

	portType := reflect.TypeOf((*port.Port)(nil)).Elem()
	nodeType := reflect.TypeOf(node).Elem()
	fieldsCount := nodeType.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := nodeType.Field(i)
		if !field.Type.ConvertibleTo(portType) {
			continue
		}
		portField := reflect.ValueOf(node).Elem().Field(i).Interface().(port.Port)
		if portField.Direction() != port.DirInput {
			continue
		}
		nodeConnectors = append(nodeConnectors, portField.Connectors()...)
	}

	return nodeConnectors
}

func (runner *ConcurrentRunner) stopNode(node Node) {
	err := node.Stop()
	if err != nil {
		panicStr := fmt.Sprintf("unable to stop node: %s %v, reason: %v", node.Name(), node.ID(), err)
		panic(panicStr)
	}
	log.Println("Stopped node:", node.Name(), node.ID())
}
