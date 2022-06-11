package command

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"

	"github.com/Workiva/go-datastructures/queue"
	"golang.org/x/crypto/blake2b"
)

// Value represents a Nyandis value.
// It can be a boolean, integer, float, byte array, string, list, map, set, or buffer.
type Value struct {
	Type        Type
	Size        int
	Value       []byte
	Values      []Value
	ValueMap    map[string]Value
	ValueSet    map[string]Value
	ValueBuffer *queue.RingBuffer
}

func (v Value) Hash() string {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte{byte(v.Type)})

	a := [8]byte{}
	binary.BigEndian.PutUint64(a[:], uint64(v.Size))
	buf.Write(a[:])
	if v.Value != nil {
		buf.Write(v.Value)
	} else {
		buf.Write([]byte{1, 2, 3, 4, 255, 254, 253, 252})
	}
	for _, v := range v.Values {
		buf.Write([]byte(v.Hash()))
	}
	for k, v := range v.ValueMap {
		buf.Write([]byte(k))
		buf.Write([]byte(v.Hash()))
	}
	for k, v := range v.ValueSet {
		buf.Write([]byte(k))
		buf.Write([]byte(v.Hash()))
	}
	hashed := blake2b.Sum256(buf.Bytes())
	return hex.EncodeToString(hashed[:])
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

func NewList(values []any) Value {
	vs := make([]Value, len(values))
	for i, v := range values {
		switch a := v.(type) {
		case bool:
			vs[i] = NewBool(a)
		case int:
			vs[i] = NewInt(int64(a))
		case int64:
			vs[i] = NewInt(a)
		case float64:
			vs[i] = NewFloat(a)
		case []byte:
			vs[i] = NewByteArray(a)
		case string:
			vs[i] = NewString(a)
		case []any:
			vs[i] = NewList(a)
		case map[string]any:
			vs[i] = NewMap(a)
		}
	}
	return Value{
		Type:   TypeList,
		Size:   len(values),
		Values: vs,
	}
}

func NewMap(values map[string]any) Value {
	vs := make(map[string]Value)
	for k, v := range values {
		switch a := v.(type) {
		case bool:
			vs[k] = NewBool(a)
		case int:
			vs[k] = NewInt(int64(a))
		case int64:
			vs[k] = NewInt(a)
		case float64:
			vs[k] = NewFloat(a)
		case []byte:
			vs[k] = NewByteArray(a)
		case string:
			vs[k] = NewString(a)
		case []any:
			vs[k] = NewList(a)
		case map[string]any:
			vs[k] = NewMap(a)
		}
	}
	return Value{
		Type:     TypeMap,
		Size:     len(values),
		ValueMap: vs,
	}
}

// TODO: create set.

func NewBuffer(size int) Value {
	return Value{
		Type:        TypeBuffer,
		Size:        size,
		ValueBuffer: queue.NewRingBuffer(uint64(size)),
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
