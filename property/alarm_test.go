package property

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlarm_SetActive(t *testing.T) {
	// GIVEN
	alarm := NewAlarm("alarm", "description", SeverityMinor)

	//WHEN
	err := alarm.SetActive(true)

	//THEN
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x1}, alarm.Data)
}

func TestAlarm_Write(t *testing.T) {
	// GIVEN
	alarm := NewAlarm("alarm", "description", SeverityMinor)
	src := []byte{0x01}

	//WHEN
	written, err := alarm.Write(src)

	//THEN
	assert.Nil(t, err)
	assert.Equal(t, len(src), written)
	assert.Equal(t, src, alarm.Data)
}

func TestAlarm_WriteEmpty(t *testing.T) {
	// GIVEN
	alarm := NewAlarm("alarm", "description", SeverityMinor)
	src := []byte{}

	//WHEN
	written, err := alarm.Write(src)

	//THEN
	assert.Nil(t, err)
	assert.Equal(t, len(src), written)
	assert.Equal(t, src, alarm.Data)
}
