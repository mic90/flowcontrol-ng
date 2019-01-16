package property

import (
	"github.com/mailru/easyjson"
	"github.com/mic90/flowcontrol-ng/buffer"
)

type NamedReadWriter interface {
	buffer.TypedReadWriter
	easyjson.Marshaler
	GetName() string
	GetDescription() string
}

type Type struct {
	buffer.Type
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        Unit   `json:"unit"`
	Persistent  bool   `json:"-"`
	Flags       Flags  `json:"flags"`
}

func New(name, description string, unit Unit, persistent bool, flags Flags) *Type {
	return &Type{*buffer.New(), name, description, unit, persistent, flags}
}

func (p *Type) Flag(flagName string) (int, error) {
	return p.Flags.Flag(flagName)
}

func (p *Type) GetName() string {
	return p.Name
}

func (p *Type) GetDescription() string {
	return p.Description
}
