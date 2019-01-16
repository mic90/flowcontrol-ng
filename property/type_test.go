package property_test

import (
	"bytes"
	"github.com/mailru/easyjson/jwriter"
	. "github.com/mic90/flowcontrol-ng/property"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProperty_MarshalJSON(t *testing.T) {
	// GIVEN
	const expectedJSON = `{"name":"propA","description":"description","unit":"%","flags":{"flag":10},"data":"ZA==","size":1}`

	value := []byte{100}
	flags := *NewFlags(Flag{"flag", 10})
	prop := New("propA", "description", UnitPercent, false, flags)
	prop.Write(value)

	// WHEN
	w := &jwriter.Writer{}
	buffer := &bytes.Buffer{}
	prop.MarshalEasyJSON(w)
	_, err := w.DumpTo(buffer)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, string(buffer.Bytes()))
}

func BenchmarkProperty_MarshalJSON(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	value := []byte{100}
	flags := *NewFlags(Flag{"flag", 10})
	prop := New("propA", "description", UnitPercent, false, flags)
	prop.Write(value)

	w := &jwriter.Writer{}
	buffer := &bytes.Buffer{}

	// WHEN
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		prop.MarshalEasyJSON(w)
		_, err := w.DumpTo(buffer)
		if err != nil {
			b.Fail()
		}
	}
}
