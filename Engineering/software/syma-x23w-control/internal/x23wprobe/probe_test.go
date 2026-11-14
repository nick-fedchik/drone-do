package x23wprobe

import (
	"bytes"
	"testing"
)

func TestFirmwareVersionQueryBytes(t *testing.T) {
	got := FirmwareVersionQuery()
	want := []byte{'F', 'D', 'W', 'N', 0x20, 0x00, 0x09, 0x00, 0x00, 0x00}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestReassembleJPEGs(t *testing.T) {
	p1 := append(make([]byte, rtpHeaderSize), 0xFF, 0xD8, 0xFF, 0xE0)
	p2 := append(make([]byte, rtpHeaderSize), 0x11, 0x22, 0xFF, 0xD9)
	frames := ReassembleJPEGs([][]byte{p1, p2})
	if len(frames) != 1 {
		t.Fatalf("frames = %d, want 1", len(frames))
	}
	want := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x11, 0x22, 0xFF, 0xD9}
	if !bytes.Equal(frames[0], want) {
		t.Fatalf("frame = %x, want %x", frames[0], want)
	}
}
