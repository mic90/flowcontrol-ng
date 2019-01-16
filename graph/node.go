package graph

import (
	"context"
	"errors"
	"github.com/rs/xid"
)

//Node is the base struct used in the flow.
//Each node consists of inputs, outputs and its properties.
//The basic flow is that, when the inputs are received from connected nodes
//the node logic is executed.
//One must set all output values in the Process function to propagate its values to other nodes inputs
type Node interface {
	Name() string
	Description() string
	Version() string
	ID() xid.ID

	Setup() error
	Process(context.Context, Timer) error
	Stop() error
}

//BaseNode is the base struct which must be added by composition to the user-defined nodes
type BaseNode struct {
	name        string
	description string
	version     string
	id          xid.ID
}

func NewBaseNode(name, description, version string) *BaseNode {
	return &BaseNode{name, description, version, xid.New()}
}

//GetName returns the name of the node.
//The name should not contain spaces, be short and meaningfull
func (node *BaseNode) Name() string {
	return node.name
}

//GetDescription returns the description of the node
//The descrption could be as long as one would like.
//Provide all important informations here
func (node *BaseNode) Description() string {
	return node.description
}

//GetVersion returns version of given node
func (node *BaseNode) Version() string {
	return node.version
}

//GetID returns node id. The id will be generated automatically by graph when the node is added
func (node *BaseNode) ID() xid.ID {
	return node.id
}

//Setup this function will be triggered when the node is started.
//It will be triggred only once. Put all the initialization code here
func (node *BaseNode) Setup() error {
	return nil
}

//Process this is the main node function.
//It will be run in the loop, when this function is executed all non-optional inputs are already provided
//Put all the node logic here.
func (node *BaseNode) Process(ctx context.Context, timer Timer) error {
	return errors.New("override me")
}

//Stop this function will be called on graph stop
//It may be left as is which will do nothing
func (node *BaseNode) Stop() error {
	return nil
}
