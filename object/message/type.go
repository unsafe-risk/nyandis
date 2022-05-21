package message

type MessageType uint8

const (
	SET MessageType = iota
	GET
	DELETE
	EXISTS
)
