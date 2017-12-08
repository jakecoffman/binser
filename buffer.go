package binserializer

import (
	"io"
	"math"
)

// Buffer is a helper struct for serializing and deserializing as the caller
// does not need to externally manage where in the buffer they are currently reading
// or writing to.
type Buffer struct {
	buf []byte // the backing byte slice
	pos int    // current position in read/write
	err error
}

// Creates a new Buffer with a backing byte slice of the provided size
func NewBuffer(size int) *Buffer {
	b := &Buffer{}
	b.buf = make([]byte, size)
	return b
}

// Creates a new Buffer using the original backing slice
func NewBufferFromBytes(buf []byte) *Buffer {
	b := &Buffer{}
	b.buf = buf
	b.pos = 0
	return b
}

// Creates a new buffer from a byte slice by copying it
func NewBufferCopyFromBytes(buf []byte) *Buffer {
	b := &Buffer{}
	b.buf = make([]byte, len(buf))
	copy(b.buf, buf)
	return b
}

// Error returns any errors saved from other operations
func (b *Buffer) Error() error {
	return b.err
}

// Copy returns a copy of Buffer
func (b *Buffer) Copy() *Buffer {
	c := NewBuffer(len(b.buf))
	copy(c.buf, b.buf)
	return c
}

// Len returns the length of the backing byte slice
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Bytes returns the backing byte slice and any errors
func (b *Buffer) Bytes() ([]byte, error) {
	return b.buf, b.err
}

// Pos returns the current position cursor
func (b *Buffer) Pos() int {
	return b.pos
}

// Resets the position back to beginning of buffer
func (b *Buffer) Reset() {
	b.pos = 0
}

// GetByte decodes a little-endian byte
func (b *Buffer) GetByte() (result byte) {
	if b.err != nil {
		return
	}
	return b.GetUint8()
}

// GetBytes returns a byte slice possibly smaller than length if bytes are not
// available from the reader.
func (b *Buffer) GetBytes(length int) (result []byte) {
	if b.err != nil {
		return
	}
	if len(b.buf) < length {
		b.err = io.EOF
		return nil
	}
	value := b.buf[b.pos: b.pos+length]
	b.pos += length
	return value
}

// GetUint8 decodes a little-endian uint8 from the buffer
func (b *Buffer) GetUint8() uint8 {
	if b.err != nil {
		return 0
	}
	if b.pos+SizeUint8 > len(b.buf) {
		b.err = io.EOF
		return 0
	}
	buf := b.buf[b.pos: b.pos+SizeUint8]
	b.pos++
	return uint8(buf[0])
}

// GetUint16 decodes a little-endian uint16 from the buffer
func (b *Buffer) GetUint16() (n uint16) {
	if b.err != nil {
		return
	}
	buf := b.GetBytes(SizeUint16)
	n |= uint16(buf[0])
	n |= uint16(buf[1]) << 8
	return n
}

// GetUint32 decodes a little-endian uint32 from the buffer
func (b *Buffer) GetUint32() (n uint32) {
	buf := b.GetBytes(SizeUint32)
	if b.err != nil {
		return
	}
	n |= uint32(buf[0])
	n |= uint32(buf[1]) << 8
	n |= uint32(buf[2]) << 16
	n |= uint32(buf[3]) << 24
	return
}

// GetUint64 decodes a little-endian uint64 from the buffer
func (b *Buffer) GetUint64() (n uint64) {
	buf := b.GetBytes(SizeUint64)
	if b.err != nil {
		return
	}
	n |= uint64(buf[0])
	n |= uint64(buf[1]) << 8
	n |= uint64(buf[2]) << 16
	n |= uint64(buf[3]) << 24
	n |= uint64(buf[4]) << 32
	n |= uint64(buf[5]) << 40
	n |= uint64(buf[6]) << 48
	n |= uint64(buf[7]) << 56
	return
}

// GetInt8 decodes a little-endian int8 from the buffer
func (b *Buffer) GetInt8() (int8) {
	if b.pos+1 > len(b.buf) {
		b.err = io.EOF
		return 0
	}
	buf := b.buf[b.pos: b.pos+SizeInt8]
	b.pos += 1
	return int8(buf[0])
}

// GetInt16 decodes a little-endian int16 from the buffer
func (b *Buffer) GetInt16() (n int16) {
	buf := b.GetBytes(SizeInt16)
	if b.err != nil {
		return
	}
	n |= int16(buf[0])
	n |= int16(buf[1]) << 8
	return
}

// GetInt32 decodes a little-endian int32 from the buffer
func (b *Buffer) GetInt32() (n int32) {
	buf := b.GetBytes(SizeInt32)
	if b.err != nil {
		return
	}
	n |= int32(buf[0])
	n |= int32(buf[1]) << 8
	n |= int32(buf[2]) << 16
	n |= int32(buf[3]) << 24
	return
}

// GetInt64 decodes a little-endian int64 from the buffer
func (b *Buffer) GetInt64() (n int64) {
	buf := b.GetBytes(SizeInt64)
	if b.err != nil {
		return
	}
	n |= int64(buf[0])
	n |= int64(buf[1]) << 8
	n |= int64(buf[2]) << 16
	n |= int64(buf[3]) << 24
	n |= int64(buf[4]) << 32
	n |= int64(buf[5]) << 40
	n |= int64(buf[6]) << 48
	n |= int64(buf[7]) << 56
	return
}

// ReadFloat32 decodes a little-endian float32 into the buffer.
func (b *Buffer) GetFloat32() float32 {
	buf := b.GetUint32()
	if b.err != nil {
		return 0
	}
	return math.Float32frombits(buf)
}

// ReadFloat64 decodes a little-endian float64 into the buffer.
func (b *Buffer) GetFloat64() float64 {
	buf := b.GetUint64()
	if b.err != nil {
		return 0
	}
	return math.Float64frombits(buf)
}

// WriteByte encodes a little-endian uint8 into the buffer.
func (b *Buffer) WriteByte(n byte) {
	b.buf[b.pos] = uint8(n)
	b.pos++
}

// WriteBytes encodes a little-endian byte slice into the buffer
func (b *Buffer) WriteBytes(src []byte) {
	for i := 0; i < len(src); i += 1 {
		b.WriteByte(uint8(src[i]))
	}
}

// WriteBytes encodes a little-endian byte slice into the buffer
func (b *Buffer) WriteBytesN(src []byte, length int) {
	for i := 0; i < length; i += 1 {
		b.WriteByte(uint8(src[i]))
	}
}

// WriteUint8 encodes a little-endian uint8 into the buffer.
func (b *Buffer) WriteUint8(n uint8) {
	b.buf[b.pos] = byte(n)
	b.pos++
}

// WriteUint16 encodes a little-endian uint16 into the buffer.
func (b *Buffer) WriteUint16(n uint16) {
	b.buf[b.pos] = byte(n)
	b.buf[b.pos+1] = byte(n >> 8)
	b.pos+=2
}

// WriteUint32 encodes a little-endian uint32 into the buffer.
func (b *Buffer) WriteUint32(n uint32) {
	b.buf[b.pos] = byte(n)
	b.buf[b.pos+1] = byte(n >> 8)
	b.buf[b.pos+2] = byte(n >> 16)
	b.buf[b.pos+3] = byte(n >> 24)
	b.pos+=4
}

// WriteUint64 encodes a little-endian uint64 into the buffer.
func (b *Buffer) WriteUint64(n uint64) {
	for i := uint(0); i < uint(SizeUint64); i++ {
		b.buf[b.pos] = byte(n >> (i * 8))
		b.pos++
	}
}

// WriteInt8 encodes a little-endian int8 into the buffer.
func (b *Buffer) WriteInt8(n int8) {
	b.buf[b.pos] = byte(n)
	b.pos++
}

// WriteInt16 encodes a little-endian int16 into the buffer.
func (b *Buffer) WriteInt16(n int16) {
	b.buf[b.pos] = byte(n)
	b.buf[b.pos+1] = byte(n >> 8)
	b.pos+=2
}

// WriteInt32 encodes a little-endian int32 into the buffer.
func (b *Buffer) WriteInt32(n int32) {
	b.buf[b.pos] = byte(n)
	b.buf[b.pos+1] = byte(n >> 8)
	b.buf[b.pos+2] = byte(n >> 16)
	b.buf[b.pos+3] = byte(n >> 24)
	b.pos+=4
}

// WriteInt64 encodes a little-endian int64 into the buffer.
func (b *Buffer) WriteInt64(n int64) {
	for i := uint(0); i < uint(SizeInt64); i++ {
		b.buf[b.pos] = byte(n >> (i * 8))
		b.pos++
	}
}

// WriteFloat32 encodes a little-endian float32 into the buffer.
func (b *Buffer) WriteFloat32(n float32) {
	b.WriteUint32(math.Float32bits(n))
}

// WriteFloat64 encodes a little-endian float64 into the buffer.
func (b *Buffer) WriteFloat64(n float64) {
	b.WriteUint64(math.Float64bits(n))
}
