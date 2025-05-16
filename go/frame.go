package blink

const (
	TypeCreate      = 0x01
	TypeSubscribe   = 0x02
	TypePublish     = 0x03
	TypeMessage     = 0x04
	TypeUnsubscribe = 0x05
	TypeKeyUpdate   = 0x06
	TypeRotateKey   = 0x07
)

type Frame interface {
	Type() byte
}

// Frame Variants

type CreateFrame struct {
	JWT       []byte
	TopicName string
	Flags     byte
}

func (f *CreateFrame) Type() byte { return TypeCreate }

type SubscribeFrame struct {
	JWT     []byte
	TopicID uint32
}

func (f *SubscribeFrame) Type() byte { return TypeSubscribe }

type UnsubscribeFrame struct {
	JWT     []byte
	TopicID uint32
}

func (f *UnsubscribeFrame) Type() byte { return TypeUnsubscribe }

type PublishFrame struct {
	JWT     []byte
	TopicID uint32
	Payload []byte
}

func (f *PublishFrame) Type() byte { return TypePublish }

type MessageFrame struct {
	TopicID uint32
	Payload []byte
}

func (f *MessageFrame) Type() byte { return TypeMessage }

type KeyUpdateFrame struct {
	TopicID uint32
	NewKey  []byte
}

func (f *KeyUpdateFrame) Type() byte { return TypeKeyUpdate }

type RotateKeyFrame struct {
	JWT     []byte
	TopicID uint32
	NewKey  []byte
}

func (f *RotateKeyFrame) Type() byte { return TypeRotateKey }