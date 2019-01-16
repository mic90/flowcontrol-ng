package port

import (
	"context"
	"log"
)

//TODO: Think of better names
type Connector interface {
	Wait()
	Notify()
	Trigger() error
	Qos() Qos
}

type PortConnector struct {
	fromPort  *Type
	toPort    *Type
	data      []byte
	writeID   int64
	writeDone chan bool
	qos       Qos
	ctx       context.Context
}

func NewConnector(ctx context.Context, qos Qos, fromPort, toPort *Type) *PortConnector {
	if toPort.Direction() != DirInput {
		panic("destination port must be of input type")
	}
	if fromPort.Direction() != DirOutput {
		panic("source port must be of output type")
	}
	// prepare sync channel, for qos2 there must me single value channel, which should block until reader is ready
	// for lower qos, channel is buffered to prevent write trigger loose.
	var dataWritten chan bool
	if qos == Qos2 {
		dataWritten = make(chan bool)
	} else {
		dataWritten = make(chan bool, 1)
	}
	connector := &PortConnector{fromPort, toPort, make([]byte, 0), 0, dataWritten, qos, ctx}

	// subscribe for the write events on the source port
	fromPort.Connect(connector)
	toPort.Connect(connector)

	return connector
}

func (c *PortConnector) Trigger() error {
	writeIdDiff := c.fromPort.WriteID() - c.writeID
	if writeIdDiff > 1 && c.qos == 2 {
		log.Println("WARNING! source port data was lost, missed writes: ", writeIdDiff)
	}

	c.matchSrcPortSize()
	_, err := c.fromPort.Read(c.data)
	if err != nil {
		return err
	}
	_, err = c.toPort.Write(c.data)
	if err != nil {
		return err
	}
	c.writeID = c.fromPort.WriteID()
	return nil
}

func (c *PortConnector) Wait() {
	if c.qos < Qos1 {
		return
	}
	log.Printf("%v Waiting for data to be written", c.fromPort.Name())
	// wait for new data at src port or external stop request
	select {
	case <-c.writeDone:
		return
	case <-c.ctx.Done():
		return
	}
}

func (c *PortConnector) Notify() {
	// write and forget
	if c.qos < Qos2 {
		select {
		case c.writeDone <- true:
			return
		default:
			return
		}
	}
	// write and wait for the target to read it
	if c.qos == Qos2 {
		select {
		case c.writeDone <- true:
			return
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *PortConnector) Qos() Qos {
	return c.qos
}

func (c *PortConnector) matchSrcPortSize() {
	size := c.fromPort.Size
	if len(c.data) >= size {
		return
	}
	newData := make([]byte, size)
	c.data = newData
}
