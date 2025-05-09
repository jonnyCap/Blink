package blink

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseSubscribeFrame(t *testing.T) {
	input := "SUBSCRIBE logs\n"
	reader := bufio.NewReader(strings.NewReader(input))

	frame, err := ParseFrame(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Command != "SUBSCRIBE" || frame.Topic != "logs" {
		t.Errorf("parsed frame is incorrect: %+v", frame)
	}
}

func TestParsePublishFrame(t *testing.T) {
	input := "PUBLISH updates 11\nhello world"
	reader := bufio.NewReader(strings.NewReader(input))

	frame, err := ParseFrame(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Command != "PUBLISH" || frame.Topic != "updates" || string(frame.Payload) != "hello world" {
		t.Errorf("parsed publish frame is incorrect: %+v", frame)
	}
}

func TestEncodeFrame(t *testing.T) {
	frame := &Frame{
		Command: "SUBSCRIBE",
		Topic:   "metrics",
	}

	data, err := EncodeFrame(frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "SUBSCRIBE metrics\n"
	if string(data) != expected {
		t.Errorf("expected %q but got %q", expected, string(data))
	}
}
