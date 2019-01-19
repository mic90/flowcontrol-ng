package types_test

import (
	. "github.com/mic90/flowcontrol-ng/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool_Deserialize(t *testing.T) {
	// GIVEN
	boolVal := NewBool()
	expectedBool := true

	// WHEN
	boolVal.Deserialize([]byte{0x01})
	deserialized := boolVal.Value

	// THEN
	assert.Equal(t, expectedBool, deserialized)
}

func TestBool_Serialize(t *testing.T) {
	// GIVEN
	boolVal := NewBool()
	boolVal.Set(false)
	expectedRaw := []byte{0x00}

	// WHEN
	serialized := boolVal.Serialize()

	// THEN
	assert.Equal(t, expectedRaw, serialized)
}

func BenchmarkBool_Serialize(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	boolVal := NewBool()
	boolVal.Set(true)

	for i := 0; i < b.N; i++ {
		serialized := boolVal.Serialize()
		if len(serialized) == 0 {
			b.Fail()
		}
	}
}

func BenchmarkBool_Deserialize(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	boolVal := NewBool()
	rawData := []byte{0x01}

	for i := 0; i < b.N; i++ {
		err := boolVal.Deserialize(rawData)
		if err != nil {
			b.Fail()
		}
	}
}
