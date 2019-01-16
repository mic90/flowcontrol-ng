package property

import (
	"github.com/mic90/flowcontrol-ng/buffer"
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
	if data[0] == 0 {
		p.Active = false
	} else {
		p.Active = true
	}
	return written, nil
}
