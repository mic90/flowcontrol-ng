package buffer_test

import (
	"encoding/json"
	. "github.com/mic90/flowcontrol-ng/buffer"
	"github.com/mic90/flowcontrol-ng/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer_WriteRead(t *testing.T) {
	// GIVEN
	const expectedBufferSize int = 2
	const expectedWriteID int64 = 1

	buffer := New()
	expectedData := []byte{12, 34}
	readData := make([]byte, expectedBufferSize)
	// WHEN
	_, writeErr := buffer.Write(expectedData)
	_, readErr := buffer.Read(readData)

	// THEN
	assert.Nil(t, writeErr)
	assert.Nil(t, readErr)
	assert.Equal(t, expectedBufferSize, buffer.Size)
	assert.Equal(t, expectedWriteID, buffer.WriteID())
	assert.Equal(t, expectedData, readData)
}

func TestBuffer_WriteReadTyped(t *testing.T) {
	// GIVEN
	const expectedBufferSize int = 2
	const expectedWriteID int64 = 1
	buffer := New()
	value := types.NewInt()
	value.Set(-1235)
	readValue := types.NewInt()

	// WHEN
	_, writeErr := buffer.WriteTyped(value)
	readErr := buffer.ReadTyped(readValue)

	// THEN
	assert.Nil(t, writeErr)
	assert.Nil(t, readErr)
	assert.Equal(t, expectedBufferSize, buffer.Size)
	assert.Equal(t, expectedWriteID, buffer.WriteID())
	assert.Equal(t, value.Value, readValue.Value)
}

func TestBuffer_MarshalJSON(t *testing.T) {
	// GIVEN
	buffer := New()
	buffer.Write([]byte{1, 2, 3, 4})
	expectedJSON := `{"data":"AQIDBA==","size":4}`

	// WHEN
	resultJSON, err := json.Marshal(buffer)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, string(resultJSON))
}

func BenchmarkBuffer_WriteParallel(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	buffer := New()
	expectedData := []byte{12, 34}
	buffer.Write(expectedData)

	// WHEN
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, writeErr := buffer.Write(expectedData)
			if writeErr != nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkBuffer_ReadParallel(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	const bufferSize = 2
	buffer := New()
	expectedData := []byte{12, 34}
	readData := make([]byte, bufferSize)

	buffer.Write(expectedData)
	// WHEN
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, readErr := buffer.Read(readData)
			if readErr != nil {
				b.Fail()
			}
			if readData[0] != expectedData[0] || readData[1] != expectedData[1] {
				b.Fail()
			}
		}
	})
}
