package property_test

import (
	"github.com/mic90/flowcontrol-ng/property"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer_Property(t *testing.T) {
	// GIVEN
	flags := *property.NewFlags()
	container := property.NewContainer(
		[]property.NamedReadWriter{
			property.New("alarmA", "description", property.UnitPercent, false, flags),
		},
		[]property.NamedReadWriter{
			property.New("propA", "description", property.UnitPercent, false, flags),
		})

	//WHEN
	prop := container.Property("propA")

	//THEN
	assert.Equal(t, "propA", prop.GetName())
	assert.Equal(t, "description", prop.GetDescription())
	//assert.Equal(t, property.UnitPercent, prop.Unit)
	//assert.Equal(t, false, prop.Persistent)
}

func TestContainer_AsJSON(t *testing.T) {
	// GIVEN
	const expectedJSON = `{"alarms":[{"name":"alarmA","description":"description","active":false,"severity":0,"data":"","size":0},{"name":"alarmB","description":"description","active":false,"severity":0,"data":"","size":0}],"properties":[{"name":"propA","description":"description","unit":"%","flags":{},"data":"","size":0},{"name":"propB","description":"description","unit":"%","flags":{},"data":"","size":0}]}`
	flags := *property.NewFlags()
	container := property.Container{
		Alarms: []property.NamedReadWriter{
			property.NewAlarm("alarmA", "description", property.SeverityMinor),
			property.NewAlarm("alarmB", "description", property.SeverityMinor),
		},
		Properties: []property.NamedReadWriter{
			property.New("propA", "description", property.UnitPercent, false, flags),
			property.New("propB", "description", property.UnitPercent, false, flags),
		}}

	//WHEN
	data, err := container.AsJSON()

	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, string(data))
}

func TestContainer_FromJSON(t *testing.T) {
	// GIVEN
	const sourceJSON = `{"alarms":[{"name":"alarmA","description":"description","unit":"%","flags":{},"data":"ZA==","size":1}],"properties":[{"name":"propA","description":"description","unit":"%","flags":{},"data":"ZA==","size":1}]}`
	flags := *property.NewFlags()
	diffValue := property.New("propA", "description", property.UnitPercent, false, flags)
	container := property.Container{
		Alarms: []property.NamedReadWriter{
			property.New("alarmA", "description", property.UnitPercent, false, flags),
		},
		Properties: []property.NamedReadWriter{
			diffValue,
		}}

	//WHEN
	_, err := diffValue.Write([]byte{0, 0, 0})
	if err != nil {
		t.Fail()
	}
	err = container.FromJSON([]byte(sourceJSON))

	assert.Nil(t, err)
	assert.Equal(t, []byte{100}, diffValue.Data)
	assert.Equal(t, 1, len(container.Properties))
	assert.Equal(t, 1, len(container.Alarms))
}

func BenchmarkContainer_AsJSON(b *testing.B) {
	b.ReportAllocs()

	// GIVEN
	flags := *property.NewFlags()
	container := property.Container{
		Alarms: []property.NamedReadWriter{
			property.NewAlarm("alarmA", "description", property.SeverityMinor),
			property.NewAlarm("alarmB", "descriptionA", property.SeverityMinor),
			property.NewAlarm("alarmC", "descriptionB", property.SeverityMinor),
		},
		Properties: []property.NamedReadWriter{
			property.New("propA", "description", property.UnitPercent, false, flags),
			property.New("propB", "descriptionA", property.UnitPercent, false, flags),
			property.New("propC", "descriptionB", property.UnitPercent, false, flags),
			property.New("propC", "descriptionB", property.UnitPercent, false, flags),
			property.New("propC", "descriptionB", property.UnitPercent, false, flags),
			property.New("propC", "descriptionB", property.UnitPercent, false, flags),
		}}

	// WHEN
	for i := 0; i < b.N; i++ {
		_, err := container.AsJSON()
		if err != nil {
			b.Fail()
		}
	}
}
