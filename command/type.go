package command

type Type = uint8

const (
	TypeBool = uint8(iota)
	TypeInt
	TypeFloat
	TypeBytes
	TypeString
	TypeList
	TypeMap
	TypeSet
	TypeSortedSet
	TypeBuffer
	TypeEnd
)

type Cmd = uint8

// R: command key value
// I: command key
// E: command key value...
const (
	IStart Cmd = iota
	IExists
	IIncrement
	IDecrement

	RStart
	RSet
	RGet

	EStart
	ECompareAndSwap
)
