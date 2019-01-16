package property

import "errors"

var ErrFlagNotFound = errors.New("flag not found")

type Flag struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Flags map[string]int

func NewFlags(flags ...Flag) *Flags {
	flagsData := make(Flags) //make(map[string]int)
	for _, flag := range flags {
		flagsData[flag.Name] = flag.Value
	}
	return &flagsData
}

func (f *Flags) Flag(flagName string) (int, error) {
	result, ok := (*f)[flagName]
	if !ok {
		return 0, ErrFlagNotFound
	}
	return result, nil
}
