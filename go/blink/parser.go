package blink

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func readByte(r io.Reader) (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	return b[0], err
}

func readBytes(r io.Reader, length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func ParseFrame(r io.Reader) (Frame, error) {
	typ, err := readByte(r)
	if err != nil {
		return nil, err
	}

	switch typ {
	case TypeCreate:
		jwtLen, _ := readByte(r)
		jwt, _ := readBytes(r, int(jwtLen))
		topicLen, _ := readByte(r)
		topicBytes, _ := readBytes(r, int(topicLen))
		topic := string(topicBytes)
		flags, _ := readByte(r)
		return &CreateFrame{JWT: jwt, TopicName: topic, Flags: flags}, nil
	case TypeSubscribe:
		jwtLen, _ := readByte(r)
		jwt, _ := readBytes(r, int(jwtLen))
		var topicID uint32
		binary.Read(r, binary.BigEndian, &topicID)
		return &SubscribeFrame{JWT: jwt, TopicID: topicID}, nil
	case TypeUnsubscribe:
		jwtLen, _ := readByte(r)
		jwt, _ := readBytes(r, int(jwtLen))
		var topicID uint32
		binary.Read(r, binary.BigEndian, &topicID)
		return &UnsubscribeFrame{JWT: jwt, TopicID: topicID}, nil
	case TypePublish:
		jwtLen, _ := readByte(r)
		jwt, _ := readBytes(r, int(jwtLen))
		var topicID uint32
		binary.Read(r, binary.BigEndian, &topicID)
		var payloadLen uint32
		binary.Read(r, binary.BigEndian, &payloadLen)
		payload, _ := readBytes(r, int(payloadLen))
		return &PublishFrame{JWT: jwt, TopicID: topicID, Payload: payload}, nil
	case TypeMessage:
		var topicID uint32
		binary.Read(r, binary.BigEndian, &topicID)
		var payloadLen uint32
		binary.Read(r, binary.BigEndian, &payloadLen)
		payload, _ := readBytes(r, int(payloadLen))
		return &MessageFrame{TopicID: topicID, Payload: payload}, nil
	case TypeKeyUpdate:
		var topicID uint32
		binary.Read(r, binary.BigEndian, &topicID)
		keyLen, _ := readByte(r)
		newKey, _ := readBytes(r, int(keyLen))
		return &KeyUpdateFrame{TopicID: topicID, NewKey: newKey}, nil
	case TypeRotateKey:
		jwtLen, _ := readByte(r)
		jwt, _ := readBytes(r, int(jwtLen))
		var topicID uint32
		binary.Read(r, binary.BigEndian, &topicID)
		keyLen, _ := readByte(r)
		newKey, _ := readBytes(r, int(keyLen))
		return &RotateKeyFrame{JWT: jwt, TopicID: topicID, NewKey: newKey}, nil
	default:
		return nil, fmt.Errorf("unknown frame type: 0x%x", typ)
	}
}

func EncodeFrame(w io.Writer, f Frame) error {
	if _,err := w.Write([]byte{f.Type()}); err != nil {
		return err
	}

	switch frame := f.(type) {
	case *CreateFrame:
		w.Write([]byte{byte(len(frame.JWT))})
		w.Write(frame.JWT)
		w.Write([]byte{byte(len(frame.TopicName))})
		w.Write([]byte(frame.TopicName))
		w.Write([]byte{frame.Flags})
	case *SubscribeFrame, *UnsubscribeFrame:
		jwt := frame.(interface{ GetJWT() []byte }).GetJWT()
		w.Write([]byte{byte(len(jwt))})
		w.Write(jwt)
		binary.Write(w, binary.BigEndian, frame.(interface{ GetTopicID() uint32 }).GetTopicID())
	case *PublishFrame:
		w.Write([]byte{byte(len(frame.JWT))})
		w.Write(frame.JWT)
		binary.Write(w, binary.BigEndian, frame.TopicID)
		binary.Write(w, binary.BigEndian, uint32(len(frame.Payload)))
		w.Write(frame.Payload)
	case *MessageFrame:
		binary.Write(w, binary.BigEndian, frame.TopicID)
		binary.Write(w, binary.BigEndian, uint32(len(frame.Payload)))
		w.Write(frame.Payload)
	case *KeyUpdateFrame:
		binary.Write(w, binary.BigEndian, frame.TopicID)
		w.Write([]byte{byte(len(frame.NewKey))})
		w.Write(frame.NewKey)
	case *RotateKeyFrame:
		w.Write([]byte{byte(len(frame.JWT))})
		w.Write(frame.JWT)
		binary.Write(w, binary.BigEndian, frame.TopicID)
		w.Write([]byte{byte(len(frame.NewKey))})
		w.Write(frame.NewKey)
	default:
		return errors.New("unknown frame type")
	}

	return nil
}
