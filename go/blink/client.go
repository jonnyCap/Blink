package blink

import (
	"bufio"
	"io"
)

func SendFrame(w io.Writer, f Frame) error {
	buf := bufio.NewWriter(w)
	if err := EncodeFrame(buf, f); err != nil {
		return err
	}
	return buf.Flush()
}

func ReadFrame(r io.Reader) (Frame, error) {
	return ParseFrame(bufio.NewReader(r))
}

func NewCreateFrame(jwt []byte, topicName string, flags byte) *CreateFrame {
	return &CreateFrame{JWT: jwt, TopicName: topicName, Flags: flags}
}

func NewSubscribeFrame(jwt []byte, topicID uint32) *SubscribeFrame {
	return &SubscribeFrame{JWT: jwt, TopicID: topicID}
}

func NewUnsubscribeFrame(jwt []byte, topicID uint32) *UnsubscribeFrame {
	return &UnsubscribeFrame{JWT: jwt, TopicID: topicID}
}

func NewPublishFrame(jwt []byte, topicID uint32, payload []byte) *PublishFrame {
	return &PublishFrame{JWT: jwt, TopicID: topicID, Payload: payload}
}

func NewMessageFrame(topicID uint32, payload []byte) *MessageFrame {
	return &MessageFrame{TopicID: topicID, Payload: payload}
}

func NewKeyUpdateFrame(topicID uint32, newKey []byte) *KeyUpdateFrame {
	return &KeyUpdateFrame{TopicID: topicID, NewKey: newKey}
}

func NewRotateKeyFrame(jwt []byte, topicID uint32, newKey []byte) *RotateKeyFrame {
	return &RotateKeyFrame{JWT: jwt, TopicID: topicID, NewKey: newKey}
}