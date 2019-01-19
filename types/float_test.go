package types_test

import (
	. "github.com/mic90/flowcontrol-ng/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat_Deserialize(t *testing.T) {
	// GIVEN
	floatVal := NewFloat()
	expectedFloat := float32(12451.112)

	// WHEN
	floatVal.Deserialize([]byte{0x73, 0x8c, 0x42, 0x46})
	deserialized := floatVal.Value

	// THEN
	assert.Equal(t, expectedFloat, deserialized)
}

func TestFloat_Serialize(t *testing.T) {
	// GIVEN
	floatVal := NewFloat()
	floatVal.Set(12451.112)
	expectedRaw := []byte{0x73, 0x8c, 0x42, 0x46}

	// WHEN
	serialized := floatVal.Serialize()

	// THEN
	assert.Equal(t, expectedRaw, serialized)
}

func BenchmarkFloat_Serialize(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	floatVal := NewFloat()
	floatVal.Set(12451.112)

	for i := 0; i < b.N; i++ {
		serialized := floatVal.Serialize()
		if len(serialized) == 0 {
			b.Fail()
		}
	}
}

func BenchmarkFloat_Deserialize(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	floatVal := NewFloat()
	rawData := []byte{0x73, 0x8c, 0x42, 0x46}

	for i := 0; i < b.N; i++ {
		err := floatVal.Deserialize(rawData)
		if err != nil {
			b.Fail()
		}
	}
}
