package binserializer

import (
	"testing"
	"math"
)

func TestStream(t *testing.T) {
	b := NewWritingStream(10)
	b.WriteByte('a')
	b.WriteBytesN([]byte("bcdefghij"), 9)

	if string(b.buf) != "abcdefghij" {
		t.Fatalf("error should have written 'abcdefghij' got '%s'\n", string(b.buf))
	}

}

func TestStream_Copy(t *testing.T) {
	b := NewWritingStream(10)
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

func TestStream_GetByte(t *testing.T) {
	buf := make([]byte, 1)
	buf[0] = 0xfe
	b := NewStream(buf)
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

func TestStream_GetBytes(t *testing.T) {
	buf := make([]byte, 2)
	buf[0] = 'a'
	buf[1] = 'b'
	b := NewStream(buf)

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

	b = NewStream(buf)

	val = b.GetBytes(3)
	if b.Error() == nil {
		t.Fatal("expected EOF")
	}
	if b.Pos() != 0 {
		t.Fatal("position expected 0 got ", b.Pos())
	}
}

func TestStream_GetInt8(t *testing.T) {
	var v1, v2 int8
	v1 = math.MaxInt8

	writer := NewWritingStream(SizeInt8)
	writer.Int8(&v1)
	reader := writer.Copy()

	reader.Int8(&v2)

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != math.MaxInt8 {
		t.Fatalf("expected 0xf got: %x\n", v2)
	}

	buf := make([]byte, SizeInt8)
	buf[0] = 0xff
	b := NewStream(buf)
	b.Int8(&v2)
	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if v2 != -1 {
		t.Fatalf("expected -1 got: %x\n", v2)
	}
}

func TestStream_Int16(t *testing.T) {
	var v1, v2 int16
	v1 = 0x0fff

	writer := NewWritingStream(SizeInt16)
	writer.Int16(&v1)
	reader := writer.Copy()
	reader.Int16(&v2)

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0x0fff {
		t.Fatalf("expected 0x0fff got: %x\n", v2)
	}

	buf := make([]byte, SizeInt16)
	buf[0] = 0xff
	buf[1] = 0xff
	b := NewStream(buf)
	b.Int16(&v2)
	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if v2 != -1 {
		t.Fatalf("expected -1 got: %x\n", v2)
	}
}

func TestStream_Int32(t *testing.T) {
	var v1, v2 int32
	v1 = 0x0fffffff

	writer := NewWritingStream(SizeInt32)
	writer.Int32(&v1)
	reader := writer.Copy()

	reader.Int32(&v2)
	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0x0fffffff {
		t.Fatalf("expected 0x0fffffff got: %x\n", v2)
	}

	buf := make([]byte, SizeInt32)
	buf[0] = 0xff
	buf[1] = 0xff
	buf[2] = 0xff
	buf[3] = 0xff
	b := NewStream(buf)
	b.Int32(&v2)
	if b.Error() != nil {
		t.Fatal(b.Error())
	}

	if v2 != -1 {
		t.Fatalf("expected -1 got: %x\n", v2)
	}
}

func TestStream_Int64(t *testing.T) {
	var v1, v2 int64
	v1 = 0xf3f3f3f3f3f3

	writer := NewWritingStream(SizeInt64)
	writer.Int64(&v1)
	reader := writer.Copy()

	reader.Int64(&v2)

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0xf3f3f3f3f3f3 {
		t.Fatalf("expected 0xf3f3f3f3f3f3 got: %x\n", v2)
	}
}

func TestStream_Uint8(t *testing.T) {
	writer := NewWritingStream(SizeUint8)
	var v uint8 = 0xff
	writer.Uint8(&v)
	reader := writer.Copy()
	reader.reading = true

	var v2 uint8
	reader.Uint8(&v2)

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0xff {
		t.Fatalf("expected 0xff got: %x\n", v2)
	}
}

func TestStream_Uint16(t *testing.T) {
	writer := NewWritingStream(SizeUint16)
	var v1 uint16 = 0xffff
	writer.Uint16(&v1)
	reader := writer.Copy()

	var v2 uint16
	reader.Uint16(&v2)
	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0xffff {
		t.Fatalf("expected 0xffff got: %x\n", v2)
	}
}

func TestStream_Uint32(t *testing.T) {
	writer := NewWritingStream(SizeUint32)
	var v1 uint32 = 0xffffffff
	writer.Uint32(&v1)
	reader := writer.Copy()

	var v2 uint32
	reader.Uint32(&v2)
	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0xffffffff {
		t.Fatalf("expected 0xffffffff got: %x\n", v2)
	}
}

func TestStream_Uint64(t *testing.T) {
	var v1, v2 uint64
	v1 = 0xffffffffffffffff

	writer := NewWritingStream(SizeUint64)
	writer.Uint64(&v1)
	reader := writer.Copy()

	reader.Uint64(&v2)

	if reader.Error() != nil {
		t.Fatal(reader.Error())
	}

	if v2 != 0xffffffffffffffff {
		t.Fatalf("expected 0xffffffffffffffff got: %x\n", v2)
	}
}

func TestStream_Len(t *testing.T) {
	b := NewWritingStream(10)
	b.WriteByte('a')
	b.WriteBytesN([]byte("bcdefghij"), 9)

	if b.Len() != 10 {
		t.Fatalf("expected length of 10 got: %d\n", b.Len())
	}
}

func TestStream_WriteBytes(t *testing.T) {
	w := NewWritingStream(10)
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

func TestStream_Float32(t *testing.T) {
	var v1, v2 float32
	v1 = math.MaxFloat32

	w := NewWritingStream(4)
	w.Float32(&v1)
	r := w.Copy()
	r.Float32(&v2)
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	if v2 != math.MaxFloat32 {
		t.Fatal("expected ", math.MaxFloat32, " got ", v2)
	}
}

func TestStream_GetFloat64(t *testing.T) {
	var v1, v2 float64
	v1 = math.MaxFloat64

	w := NewWritingStream(8)
	w.Float64(&v1)
	r := w.Copy()
	r.Float64(&v2)
	if r.Error() != nil {
		t.Fatal(r.Error())
	}

	if v2 != math.MaxFloat64 {
		t.Fatal("expected ", math.MaxFloat64, " got ", v2)
	}
}
