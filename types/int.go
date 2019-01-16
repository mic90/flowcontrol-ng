package types

import (
	"errors"
	"unsafe"
)

const intUnitSize = 2
const intUnitName = "int"

type Int struct {
	TypeBase
	Value    int16 `json:"value"`
	oldValue int16
}

func NewInt() *Int {
	base := TypeBase{intUnitName, intUnitSize}
	return &Int{base, 0, 0}
}

func (t *Int) Set(value int16) {
	t.oldValue = t.Value
	t.Value = value
}

func (t *Int) Serialize() []byte {
	return (*[intUnitSize]byte)(unsafe.Pointer(&t.Value))[:]
}

func (t *Int) Deserialize(data []byte) error {
	expectedSize := intUnitSize
	if len(data) < expectedSize {
		return errors.New("array index out of bounds")
	}
	t.oldValue = t.Value

	sliced := data[:expectedSize]
	value := *(*int16)(unsafe.Pointer(&sliced[0]))
	t.Value = value

	return nil
}

func (t *Int) HasChanged() bool {
	return t.oldValue != t.Value
}
