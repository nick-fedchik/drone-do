package x23wprobe

import (
	"bytes"
	"image"
	"image/jpeg"
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

func TestFullJPEGDecodeRejectsTruncatedFrame(t *testing.T) {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatal(err)
	}
	valid := buf.Bytes()
	if _, err := jpeg.Decode(bytes.NewReader(valid)); err != nil {
		t.Fatalf("valid JPEG failed decode: %v", err)
	}
	truncated := append(append([]byte(nil), valid[:len(valid)/2]...), 0xFF, 0xD9)
	if _, err := jpeg.Decode(bytes.NewReader(truncated)); err == nil {
		t.Fatal("truncated JPEG unexpectedly decoded")
	}
}
