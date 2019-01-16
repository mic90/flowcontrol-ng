package types

type Serializer interface {
	Serialize() []byte
}

type Deserializer interface {
	Deserialize([]byte) error
}

type Describer interface {
	Name() string
	UnitSize() int
}

type Observer interface {
	HasChanged() bool
}

type Type interface {
	Serializer
	Deserializer
	Describer
	Observer
}

type TypeBase struct {
	name string
	size int
}

func NewBase(name string, size int) *TypeBase {
	return &TypeBase{name, size}
}

func (t *TypeBase) UnitSize() int {
	return t.size
}

func (t *TypeBase) Name() string {
	return t.name
}
