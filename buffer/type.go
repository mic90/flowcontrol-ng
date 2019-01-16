package buffer

import (
	"errors"
	"github.com/mic90/flowcontrol-ng/types"
	"io"
	"sync"
)

type TypedReader interface {
	io.Reader
	ReadTyped(dst types.Deserializer) error
	GetSize() int
}

type TypedWriter interface {
	io.Writer
	WriteTyped(src types.Serializer) (int, error)
}

type TypedReadWriter interface {
	TypedReader
	TypedWriter
}

// ErrNotEnoughSpace is returned when destination buffer is no table to fit written Data
var ErrNotEnoughSpace = errors.New("buffer could not fit required Data")

const initialSize = 0

// Type represents Data type which is internally stored as bytes array
// It will reallocate new memory if the written Data would not fit in the internal array
// All write/read methods are thread safe
// Buffer keeps track of the writes counts made, accessible by the WriteID method
type Type struct {
	Data    []byte `json:"data"`
	Size    int    `json:"size"`
	writeID int64
	mutex   sync.RWMutex
}

// New returns new buffer object
func New() *Type {
	data := make([]byte, initialSize)
	return &Type{data, initialSize, 0, sync.RWMutex{}}
}

func (b *Type) GetSize() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.Size
}

// WriteID returns number of writes made to the buffer
func (b *Type) WriteID() int64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.writeID
}

// WriteTyped writes Data hold by custom type which must implement Serializer
func (b *Type) WriteTyped(src types.Serializer) (int, error) {
	return b.Write(src.Serialize())
}

// ReadTypeds reads internal buffer Data as custom type which must implement Deserializer
func (b *Type) ReadTyped(dst types.Deserializer) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return dst.Deserialize(b.Data)
}

// Write writes raw bytes Data into buffer
func (b *Type) Write(data []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	srcLen := len(data)
	b.equalizeSize(srcLen)
	copy(b.Data, data)

	b.writeID++
	return srcLen, nil
}

// Read read raw internal Data into destination bytes slice
func (b *Type) Read(data []byte) (int, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	dstLen := len(data)
	if dstLen == 0 {
		return 0, nil
	}
	if dstLen < b.Size {
		return 0, ErrNotEnoughSpace
	}

	return copy(data, b.Data), nil
}

func (b *Type) equalizeSize(srcLen int) {
	if srcLen == b.Size {
		return
	}
	newData := make([]byte, srcLen)
	b.Data = newData
	b.Size = srcLen
}
