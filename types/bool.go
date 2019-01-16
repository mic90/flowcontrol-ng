package types

import (
	"errors"
	"unsafe"
)

const boolUnitSize = 1
const boolUnitName = "int"

type Bool struct {
	TypeBase
	Value    bool
	oldValue bool
}

func NewBool() *Bool {
	base := TypeBase{boolUnitName, boolUnitSize}
	return &Bool{base, false, false}
}

func (t *Bool) Set(value bool) {
	t.oldValue = t.Value
	t.Value = value
}

func (t *Bool) Serialize() []byte {
	return (*[boolUnitSize]byte)(unsafe.Pointer(&t.Value))[:]
}

func (t *Bool) Deserialize(data []byte) error {
	expectedSize := boolUnitSize
	if len(data) < expectedSize {
		return errors.New("array index out of bounds")
	}
	t.oldValue = t.Value

	sliced := data[:expectedSize]
	value := *(*bool)(unsafe.Pointer(&sliced[0]))
	t.Value = value

	return nil
}

func (t *Bool) HasChanged() bool {
	return t.oldValue != t.Value
}
