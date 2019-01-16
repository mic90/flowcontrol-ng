package port

// Qos represents required QOS for the graph connections
type Qos int

const (
	// Qos0 is the basic QOS requirement. When used, there is no guarantee
	// that the data written to node output port will be delivered to
	// the other node input port. It will not block the reader until its data is read nor
	// the writer until the reader is ready for accepting it.
	Qos0 Qos = iota
	// Qos1 is the middle level QOS requirement. When used, it will block the reader
	// until the data to be read is ready. The writer will not be blocked, so there is still
	// chance that some writes will be missed by the reader
	Qos1
	// QOs2 is the highest QOS requirement. Whe used, its guaranteed that the data written
	// to the port will be delivered to the reader. Write operation will block the writer loop
	// until data was read from the connection
	Qos2
)
