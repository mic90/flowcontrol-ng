package graph

import (
	"context"
	"github.com/mic90/flowcontrol-ng/port"
	"github.com/mic90/flowcontrol-ng/property"
	"github.com/rs/xid"
	"log"
)

type Type struct {
	runner     Runner
	nodes      map[xid.ID]Node
	connectors []port.Connector
	properties map[string]*property.Type
	ctx        context.Context
}

func New(ctx context.Context, runner Runner) *Type {
	nodes := make(map[xid.ID]Node)
	connectors := make([]port.Connector, 0, 1)
	properties := make(map[string]*property.Type)

	return &Type{runner, nodes, connectors, properties, ctx}
}

func (graph *Type) AddNode(node Node) {
	log.Printf("Adding node: %s | %s | %s\n", node.Name(), node.Version(), node.ID())

	graph.nodes[node.ID()] = node
}

func (graph *Type) AddConnector(fromPort, toPort *port.Type, qos port.Qos) {
	log.Printf("Adding connector from %s to %s", fromPort.Name(), toPort.Name())

	connector := port.NewConnector(graph.ctx, qos, fromPort, toPort)
	graph.connectors = append(graph.connectors, connector)
}

func (graph *Type) Run() {
	graph.runner.Run(graph.ctx, graph.nodes)
}

func (graph *Type) Stop() {
	graph.runner.Stop()
}
