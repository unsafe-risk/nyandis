package command

type Command struct {
	Command uint8
	Key     string
	Value   Value
	Values  []Value
}
