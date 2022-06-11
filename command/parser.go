package command

import (
	"bufio"
	"io"
)

func Parse(reader io.Reader) (Command, error) {
	cmd := Command{}

	scanner := bufio.NewReader(reader)

	// Read command type
	b, err := scanner.ReadByte()
	if err != nil {
		return cmd, err
	}
	cmd.Command = uint8(b)

	// Read key string
	key, err := scanner.ReadString(' ')
	if err != nil && err != io.EOF {
		return cmd, err
	}
	cmd.Key = key

	if cmd.Command < RStart {
		return cmd, nil
	}

	// Read value
	value, err := scanner.ReadSlice(' ')
	if err != nil && err != io.EOF {
		return cmd, err
	}
	typ, err := TypeFromBytes(value)
	if err != nil {
		return cmd, err
	}
	cmd.Value = Value{
		Type:  typ,
		Value: value[1:],
	}

	if cmd.Command < EStart {
		return cmd, nil
	}

	// Read all values
	for {
		value, err := scanner.ReadSlice(' ')
		if err != nil && err != io.EOF {
			return cmd, err
		}
		typ, err2 := TypeFromBytes(value)
		if err2 != nil {
			return cmd, err2
		}
		cmd.Values = append(cmd.Values, Value{
			Type:  typ,
			Value: value[1:],
		})
		if err == io.EOF {
			break
		}
	}

	return cmd, nil
}
