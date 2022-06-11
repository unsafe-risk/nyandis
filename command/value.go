package command

import (
	"encoding/binary"
	"errors"
	"math"
)

// Value represents a Nyandis value.
// It can be a boolean, integer, float, byte array, string, list, map, set, sorted set, or buffer.
type Value struct {
	Type  Type
	Size  int
	Value []byte
}

func NewBool(value bool) Value {
	b := byte(0)
	if value {
		b = 1
	}
	return Value{
		Type:  TypeBool,
		Size:  1,
		Value: []byte{b},
	}
}

func NewInt(value int64) Value {
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], uint64(value))
	return Value{
		Type:  TypeInt,
		Size:  8,
		Value: b[:],
	}
}

func NewFloat(value float64) Value {
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], math.Float64bits(value))
	return Value{
		Type:  TypeFloat,
		Size:  8,
		Value: b[:],
	}
}

func NewByteArray(value []byte) Value {
	return Value{
		Type:  TypeBytes,
		Size:  len(value),
		Value: value,
	}
}

func NewString(value string) Value {
	return Value{
		Type:  TypeString,
		Size:  len(value),
		Value: []byte(value),
	}
}

func ValueToBytes(value Value) ([]byte, error) {
	if value.Type >= TypeEnd {
		return nil, errors.New("invalid type")
	}
	b := make([]byte, 9+value.Size)
	b[0] = byte(value.Type)
	binary.BigEndian.PutUint64(b[1:], uint64(value.Size))
	copy(b[9:], value.Value)
	return b, nil
}

func BytesToValue(data []byte) (Value, error) {
	typ, err := TypeFromBytes(data)
	if err != nil {
		return Value{}, err
	}
	size := int64(binary.BigEndian.Uint64(data[1:9]))
	if int64(len(data)) < 9+size {
		return Value{}, errors.New("invalid value")
	}
	return Value{
		Type:  typ,
		Size:  int(size),
		Value: data[9 : 9+size],
	}, nil
}

func TypeFromBytes(data []byte) (Type, error) {
	if len(data) < 1 {
		return 0, errors.New("not enough data")
	}
	typ := Type(data[0])
	if typ >= TypeEnd {
		return 0, errors.New("invalid type")
	}
	return typ, nil
}
