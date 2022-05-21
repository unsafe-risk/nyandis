package object

type ObjectType uint8

const (
	BOOLEAN ObjectType = iota
	INTEGER
	FLOAT
	BYTES
	STRING
	ARRAY
	MAP
	SET
	BITSET
)
