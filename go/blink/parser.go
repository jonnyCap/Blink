// blink/parser.go
package blink

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseFrame(reader *bufio.Reader) (*Frame, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		return nil, errors.New("invalid frame header")
	}

	frame := &Frame{
		Command: parts[0],
		Topic:   parts[1],
		Length:  0,
		Payload: nil,
	}

	if frame.Command == "PUBLISH" {
		if len(parts) != 3 {
			return nil, errors.New("missing payload length")
		}
		length, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, errors.New("invalid payload length")
		}
		frame.Length = length
		frame.Payload = make([]byte, length)
		_, err = io.ReadFull(reader, frame.Payload)
		if err != nil {
			return nil, err
		}
	}

	return frame, nil
}

func EncodeFrame(f *Frame) ([]byte, error) {
	switch f.Command {
	case "SUBSCRIBE", "AUTH":
		return []byte(fmt.Sprintf("%s %s\n", f.Command, f.Topic)), nil
	case "PUBLISH":
		return []byte(fmt.Sprintf("%s %s %d\n%s", f.Command, f.Topic, len(f.Payload), f.Payload)), nil
	default:
		return nil, errors.New("unknown command")
	}
}