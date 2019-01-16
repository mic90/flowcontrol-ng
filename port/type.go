package port

import (
	"github.com/mic90/flowcontrol-ng/buffer"
	"github.com/mic90/flowcontrol-ng/types"
)

type Port interface {
	Name() string
	Description() string
	Direction() Dir
	Connectors() []Connector
}

type Type struct {
	buffer.Type
	name        string
	description string
	direction   Dir
	connectors  []Connector
}

func New(name, description string, direction Dir) *Type {
	return &Type{*buffer.New(), name, description, direction, make([]Connector, 0, 1)}
}

func (p *Type) Name() string {
	return p.name
}

func (p *Type) Description() string {
	return p.description
}

func (p *Type) Direction() Dir {
	return p.direction
}

func (p *Type) Connectors() []Connector {
	return p.connectors
}

func (p *Type) WriteTyped(src types.Serializer) (int, error) {
	return p.Write(src.Serialize())
}

func (p *Type) Write(data []byte) (int, error) {
	written, err := p.Type.Write(data)
	if err != nil {
		return written, err
	}
	// notify connectors if this is the output port
	if p.direction == DirOutput {
		//TODO notify ports asynchronously, as this might block
		for i := 0; i < len(p.connectors); i++ {
			conn := p.connectors[i]
			conn.Notify()
		}
	}
	return written, err
}

func (p *Type) Connect(connector Connector) {
	if p.direction == DirInput && len(p.connectors) != 0 {
		panic("input port can have only one connector attached")
	}
	p.connectors = append(p.connectors, connector)
}
