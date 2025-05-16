package blink

import (
	"bytes"
	"testing"
)

func TestPublishFrameEncodedDecode(t *testing.T){
	original := NewPublishFrame([]byte("jwt123"), 43, []byte("hello world"))

	buf := new(bytes.Buffer)
	err := EncodeFrame(buf, original)
	if err != nil {
		t.Fatalf("ecode failed: %v", err)
	}

	decoded, err := ParseFrame(buf)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	pub, ok := decoded.(*PublishFrame)
	if !ok {
		t.Fatalf("decoded frame has wrong type: %T", decoded)
	}

	if pub.TopicID != original.TopicID {
		t.Errorf("TopicID mismatch: got %v, want %v", pub.TopicID, original.TopicID)
	}
	if string(pub.JWT) != string(original.JWT) {
		t.Errorf("JWT mismatch: got %q, want %q", pub.JWT, original.JWT)
	}
	if string(pub.Payload) != string(original.Payload) {
		t.Errorf("Payload mismatch: got %q, want %q", pub.Payload, original.Payload)
	}
}