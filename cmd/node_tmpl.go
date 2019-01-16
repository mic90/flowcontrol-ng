package main

const NodeVersion = "1.0.0"

const NodeTemplate = `package {{.PackageName}}

import (
	"context"
	"github.com/mic90/flowcontrol-ng/graph"
	"github.com/mic90/flowcontrol-ng/port"
	"github.com/mic90/flowcontrol-ng/types"
	"log"
)

type {{.Component}} struct {
	graph.BaseNode

	//ports
	Input  *port.Type
	Output *port.Type

	// internal state
	value  *types.Int
}

func New{{.Component}}() *{{.Component}} {
	name := "{{.Component}}"
	description := "This is node description"
	version := "{{.Version}}"
	node := {{.Component}}{BaseNode: *graph.NewBaseNode(name, description, version)}

	node.value = types.NewInt()

	node.Input = port.New("Input", "node input port", port.DirInput)
	node.Output = port.New("Output", "node output port", port.DirOutput)

	return &node
}

func (node *{{.Component}}) Process(ctx context.Context, timer graph.Timer) error {
	node.Input.ReadTyped(node.value)
	if node.value.HasChanged() {
		log.Printf("Read input: %v", node.value.Value())

		node.Output.WriteTyped(node.value)
	}
	return nil
}`
