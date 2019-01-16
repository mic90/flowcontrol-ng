package graph

type State byte

const (
	StateStopped State = iota
	StateStarting
	StateStarted
	StateStopping
	StateError
)
