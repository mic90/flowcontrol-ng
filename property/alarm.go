package property

import (
	"github.com/mic90/flowcontrol-ng/buffer"
	"github.com/mic90/flowcontrol-ng/types"
)

type Alarm struct {
	buffer.Type
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	Severity    Severity `json:"severity"`
}

func NewAlarm(name, description string, severity Severity) *Alarm {
	return &Alarm{*buffer.New(), name, description, false, severity}
}

func (p *Alarm) GetName() string {
	return p.Name
}

func (p *Alarm) GetDescription() string {
	return p.Description
}

func (p *Alarm) SetActive(state bool) error {
	p.Active = state
	data := []byte{0x0}
	if state == true {
		data[0] = 0x01
	}
	_, err := p.Type.Write(data)
	return err
}

func (p *Alarm) Write(data []byte) (int, error) {
	written, err := p.Type.Write(data)
	if err != nil {
		return written, err
	}
	if len(data) < 1 {
		return written, nil
	}
	p.Active = bytesToBool(data)
	return written, nil
}

func (p *Alarm) WriteTyped(src types.Serializer) (int, error) {
	written, err := p.Type.WriteTyped(src)
	if err != nil {
		return written, err
	}
	bytes := src.Serialize()
	p.Active = bytesToBool(bytes)
	return written, nil
}

func bytesToBool(data []byte) bool {
	if data[0] == 0 {
		return false
	}
	return true
}