package binserializer

import (
	"testing"
	"math"
)

func TestBuffer(t *testing.T) {
	b := NewBuffer(10)
	b.WriteByte('a')
	b.WriteBytesN([]byte("bcdefghij"), 9)

	if string(b.buf) != "abcdefghij" {
		t.Fatalf("error should have written 'abcdefghij' got '%s'\n", string(b.buf))
	}

}

func TestBuffer_Copy(t *testing.T) {
	b := NewBuffer(10)
	b.WriteByte('a')
	b.WriteBytesN([]byte("bcdefghij"), 9)

	r := b.Copy()
	if r.Len() != b.Len() {
		t.Fatalf("expected copy length to be same got: %d and %d\n", r.Len(), b.Len())
	}

	data := r.GetBytes(10)
	if r.Error() != nil {
		t.Fatalf("error reading bytes from copy: %s\n", r.Error())
	}

	if string(data) != "abcdefghij" {
		t.Fatalf("error expeced: %s got %d\n", "abcdefghij", string(data))
	}

}

func TestBuffer_GetByte(t *testing.T) {
	buf := make([]byte, 1)
	buf[0] = 0xfe
	b := NewBufferFromBytes(buf)
	val := b.GetByte()

	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if val != 0xfe {
		t.Fatalf("expected 0xfe got: %x\n", val)
	}

	if b.Pos() != 1 {
		t.Fatal("position expected 1 got", b.pos)
	}
}

func TestBuffer_GetBytes(t *testing.T) {
	buf := make([]byte, 2)
	buf[0] = 'a'
	buf[1] = 'b'
	b := NewBufferFromBytes(buf)

	val := b.GetBytes(2)

	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if string(val) != "ab" {
		t.Fatalf("expected ab got: %s\n", val)
	}
	if b.Pos() != 2 {
		t.Fatal("Expected 2 got", b.Pos())
	}

	b = NewBufferFromBytes(buf)

	val = b.GetBytes(3)
	if b.Error() == nil {
		t.Fatal("expected EOF")
	}
	if b.Pos() != 0 {
		t.Fatal("position expected 0 got ", b.Pos())
	}
}

func TestBuffer_GetInt8(t *testing.T) {
	writer := NewBuffer(SizeInt8)
	writer.WriteInt8(0x0f)
	reader := writer.Copy()

	val := reader.GetInt8()

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0xf {
		t.Fatalf("expected 0xf got: %x\n", val)
	}

	buf := make([]byte, SizeInt8)
	buf[0] = 0xff
	b := NewBufferFromBytes(buf)
	val = b.GetInt8()
	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if val != -1 {
		t.Fatalf("expected -1 got: %x\n", val)
	}
}

func TestBuffer_GetInt16(t *testing.T) {
	writer := NewBuffer(SizeInt16)
	writer.WriteInt16(0x0fff)
	reader := writer.Copy()
	val := reader.GetInt16()

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0x0fff {
		t.Fatalf("expected 0x0fff got: %x\n", val)
	}

	buf := make([]byte, SizeInt16)
	buf[0] = 0xff
	buf[1] = 0xff
	b := NewBufferFromBytes(buf)
	val = b.GetInt16()
	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if val != -1 {
		t.Fatalf("expected -1 got: %x\n", val)
	}
}

func TestBuffer_GetInt32(t *testing.T) {
	writer := NewBuffer(SizeInt32)
	writer.WriteInt32(0x0fffffff)
	reader := writer.Copy()

	val := reader.GetInt32()
	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0x0fffffff {
		t.Fatalf("expected 0x0fffffff got: %x\n", val)
	}

	buf := make([]byte, SizeInt32)
	buf[0] = 0xff
	buf[1] = 0xff
	buf[2] = 0xff
	buf[3] = 0xff
	b := NewBufferFromBytes(buf)
	val = b.GetInt32()
	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if val != -1 {
		t.Fatalf("expected -1 got: %x\n", val)
	}
}

func TestBuffer_GetInt64(t *testing.T) {
	writer := NewBuffer(SizeInt64)
	writer.WriteInt64(0xf3f3f3f3f3f3)
	reader := writer.Copy()

	val := reader.GetInt64()

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0xf3f3f3f3f3f3 {
		t.Fatalf("expected 0xf3f3f3f3f3f3 got: %x\n", val)
	}
}

func TestBuffer_GetUint8(t *testing.T) {
	writer := NewBuffer(SizeUint8)
	writer.WriteUint8(0xff)
	reader := writer.Copy()

	val := reader.GetUint8()

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0xff {
		t.Fatalf("expected 0xff got: %x\n", val)
	}
}

func TestBuffer_GetUint16(t *testing.T) {
	writer := NewBuffer(SizeUint16)
	writer.WriteUint16(0xffff)
	reader := writer.Copy()

	val := reader.GetUint16()
	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0xffff {
		t.Fatalf("expected 0xffff got: %x\n", val)
	}

}

func TestBuffer_GetUint32(t *testing.T) {
	writer := NewBuffer(SizeUint32)
	writer.WriteUint32(0xffffffff)
	reader := writer.Copy()

	val := reader.GetUint32()
	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0xffffffff {
		t.Fatalf("expected 0xffffffff got: %x\n", val)
	}
}

func TestBuffer_GetUint64(t *testing.T) {
	writer := NewBuffer(SizeUint64)
	writer.WriteUint64(0xffffffffffffffff)
	reader := writer.Copy()

	val := reader.GetUint64()

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if val != 0xffffffffffffffff {
		t.Fatalf("expected 0xffffffffffffffff got: %x\n", val)
	}
}

func TestBuffer_Len(t *testing.T) {
	b := NewBuffer(10)
	b.WriteByte('a')
	b.WriteBytesN([]byte("bcdefghij"), 9)

	if b.Len() != 10 {
		t.Fatalf("expected length of 10 got: %d\n", b.Len())
	}
}

func TestBuffer_WriteBytes(t *testing.T) {
	w := NewBuffer(10)
	w.WriteBytes([]byte("0123456789"))
	r := w.Copy()
	val := r.GetBytes(10)
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	if string(val) != "0123456789" {
		t.Fatalf("expected 0123456789 got: %s %d\n", val, len(val))
	}

}

func TestBuffer_GetFloat32(t *testing.T) {
	w := NewBuffer(4)
	w.WriteFloat32(math.MaxFloat32)
	r := w.Copy()
	result := r.GetFloat32()
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	if result != math.MaxFloat32 {
		t.Fatal("expected ", math.MaxFloat32, " got ", result)
	}
}

func TestBuffer_GetFloat64(t *testing.T) {
	w := NewBuffer(8)
	w.WriteFloat64(math.MaxFloat64)
	r := w.Copy()
	result := r.GetFloat64()
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	if result != math.MaxFloat64 {
		t.Fatal("expected ", math.MaxFloat32, " got ", result)
	}
}
