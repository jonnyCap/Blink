package blink

type Frame struct {
	Command string // AUTH, SUBSCRIBE, PUBLISH
	Topic   string
	Length  int
	Payload []byte
}