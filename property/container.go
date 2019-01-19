package property

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson/jwriter"
)

type Container struct {
	Alarms        []NamedReadWriter `json:"alarms"`
	Properties    []NamedReadWriter `json:"properties"`
	jsonWriter    jwriter.Writer
	buffer        bytes.Buffer
	alarmsMap     map[string]int
	propertiesMap map[string]int
}

func NewContainer(alarms []NamedReadWriter, properties []NamedReadWriter) *Container {
	propertiesMap := make(map[string]int)
	alarmsMap := make(map[string]int)
	for i := range alarms {
		alarmsMap[alarms[i].GetName()] = i
	}
	for i := range properties {
		propertiesMap[properties[i].GetName()] = i
	}
	return &Container{Alarms: alarms, Properties: properties, alarmsMap: alarmsMap, propertiesMap: propertiesMap}
}

func (c *Container) Property(name string) NamedReadWriter {
	index, ok := c.propertiesMap[name]
	if !ok {
		panic(fmt.Sprintf("Unable to find property named %s", name))
	}
	return c.Properties[index]
}

func (c *Container) Alarm(name string) NamedReadWriter {
	index, ok := c.alarmsMap[name]
	if !ok {
		panic(fmt.Sprintf("Unable to find alarm named %s", name))
	}
	return c.Alarms[index]
}

func (c *Container) AsJSON() ([]byte, error) {
	c.buffer.Reset()
	c.buffer.WriteString(`{"alarms":[`)
	alarmsCount := len(c.Alarms)
	for i := range c.Alarms {
		c.Alarms[i].MarshalEasyJSON(&c.jsonWriter)
		_, err := c.jsonWriter.DumpTo(&c.buffer)
		if err != nil {
			return []byte{}, err
		}
		if i < alarmsCount - 1 {
			c.buffer.WriteByte(',')
		}
	}
	c.buffer.WriteString(`],`)
	c.buffer.WriteString(`"properties":[`)
	propertiesCount := len(c.Properties)
	for i := range c.Properties {
		c.Properties[i].MarshalEasyJSON(&c.jsonWriter)
		_, err := c.jsonWriter.DumpTo(&c.buffer)
		if err != nil {
			return []byte{}, err
		}
		if i < propertiesCount - 1 {
			c.buffer.WriteByte(',')
		}
	}
	c.buffer.WriteString(`]}`)
	return c.buffer.Bytes(), nil
}

func (c *Container) FromJSON(data []byte) error {
	var container map[string][]struct {
		Name string
		Data []byte
		Size int
	}
	//var container Container
	err := json.Unmarshal(data, &container)
	if err != nil {
		return err
	}
	properties, ok := container["properties"]
	if !ok {
		return fmt.Errorf("unable to find properties value in provided data")
	}
	for i := range properties {
		loadedProp := properties[i]

		for innerI := range c.Properties {
			prop := c.Properties[innerI]
			if prop.GetName() == loadedProp.Name {
				written, err := prop.Write(loadedProp.Data)
				if err != nil {
					return err
				}
				if written != loadedProp.Size {
					return fmt.Errorf("write failed, wrriten: %d, expected to write: %d", written, loadedProp.Size)
				}
			}
		}
	}
	return nil
}
