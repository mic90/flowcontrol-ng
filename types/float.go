package types

import (
	"errors"
	"unsafe"
)

const floatUnitSize = 4
const floatUnitName = "float"

type Float struct {
	TypeBase
	Value    float32
	oldValue float32
}

func NewFloat() *Float {
	return &Float{TypeBase{floatUnitName, floatUnitSize}, 0.0, 0.0}
}

func (t *Float) Set(value float32) {
	t.oldValue = t.Value
	t.Value = value
}

func (t *Float) Serialize() []byte {
	return (*[floatUnitSize]byte)(unsafe.Pointer(&t.Value))[:]
}

func (t *Float) Deserialize(data []byte) error {
	expectedSize := floatUnitSize
	if len(data) < expectedSize {
		return errors.New("array index out of bounds")
	}
	sliced := data[:expectedSize]
	value := *(*float32)(unsafe.Pointer(&sliced[0]))
	t.Value = value

	return nil
}
