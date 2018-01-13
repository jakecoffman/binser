package binserializer

import (
	"io"
	"math"
)

// Stream is a helper struct for serializing and deserializing as the caller
// does not need to externally manage where in the buffer they are currently reading
// or writing to.
type Stream struct {
	buf     []byte // the backing byte slice
	pos     int    // current position in read/write: TODO: use internal slice pos instead
	reading bool   // is the buffer for reading or writing?
	err     error  // records errors reading or writing
}

// Creates a new writing Stream with a backing byte slice of the provided size
func NewWritingStream(size int) Stream {
	return Stream{
		buf: make([]byte, 0, size),
	}
}

// Creates a new Stream using the original backing slice
// If a buffer is provided with a length, then it will be a read-only stream
// If a buffer has no length but a capacity, then it will be a write-only stream
// If nil is provided then it will be a write-only stream with a new buffer allocated
func NewStream(buf []byte) Stream {
	if len(buf) == 0 {
		if cap(buf) > 0 {
			return Stream{buf: buf}
		}
		// max MTU size seems like a good default
		return NewWritingStream(1500)
	}
	return Stream {
		buf:     buf,
		pos:     0,
		reading: true,
	}
}

// Creates a new buffer from a byte slice by copying it
func NewReadingStreamCopy(buf []byte) Stream {
	b := make([]byte, len(buf))
	copy(b, buf)
	return Stream{
		buf: b,
		reading: true,
	}
}

// Error returns any errors saved from other operations
func (b *Stream) Error() error {
	return b.err
}

// IsReading returns true if the stream is read-only, false if write-only
func (b Stream) IsReading() bool {
	return b.reading
}

// Copy returns a copy of the Stream in read-only mode
func (b Stream) Copy() Stream {
	return NewReadingStreamCopy(b.buf)
}

// Len returns the length of the backing byte slice
func (b *Stream) Len() int {
	return len(b.buf)
}

// Bytes returns the backing byte slice and any errors
func (b *Stream) Bytes() ([]byte, error) {
	return b.buf, b.err
}

// Pos returns the current position cursor
func (b *Stream) Pos() int {
	return b.pos
}

// Resets the position back to beginning of buffer
func (b *Stream) Reset() {
	b.pos = 0
}

// GetBytes returns a byte slice possibly smaller than length if bytes are not
// available from the reader.
func (b *Stream) GetBytes(length int) (result []byte) {
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

// Uint8 reads or writes a uint8
func (b *Stream) Uint8(n *uint8) {
	if b.err != nil {
		return
	}
	if b.reading {
		if b.pos+SizeUint8 > len(b.buf) {
			b.err = io.EOF
			return
		}
		buf := b.buf[b.pos: b.pos+SizeUint8]
		b.pos++
		*n = uint8(buf[0])
		return
	}
	b.buf = append(b.buf, *n)
	b.pos++
	return
}

// Uint16 reads or writes a uint16
func (b *Stream) Uint16(n *uint16) {
	if b.err != nil {
		return
	}
	if b.reading {
		buf := b.GetBytes(SizeUint16)
		var v uint16
		v |= uint16(buf[0])
		v |= uint16(buf[1]) << 8
		*n = v
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.buf = append(b.buf, byte(*n >> 8))
	b.pos += 2
	return
}

// Uint32 reads or writes a uint32
func (b *Stream) Uint32(n *uint32) {
	if b.err != nil {
		return
	}
	if b.reading {
		buf := b.GetBytes(SizeUint32)
		var v uint32
		v |= uint32(buf[0])
		v |= uint32(buf[1]) << 8
		v |= uint32(buf[2]) << 16
		v |= uint32(buf[3]) << 24
		*n = v
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.buf = append(b.buf, byte(*n >> 8))
	b.buf = append(b.buf, byte(*n >> 16))
	b.buf = append(b.buf, byte(*n >> 24))
	b.pos += 4
	return
}

// Uint64 reads or writes a uint64
func (b *Stream) Uint64(n *uint64) {
	if b.err != nil {
		return
	}
	if b.reading {
		buf := b.GetBytes(SizeUint64)
		var v uint64
		v |= uint64(buf[0])
		v |= uint64(buf[1]) << 8
		v |= uint64(buf[2]) << 16
		v |= uint64(buf[3]) << 24
		v |= uint64(buf[4]) << 32
		v |= uint64(buf[5]) << 40
		v |= uint64(buf[6]) << 48
		v |= uint64(buf[7]) << 56
		*n = v
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.buf = append(b.buf, byte(*n >> 8))
	b.buf = append(b.buf, byte(*n >> 16))
	b.buf = append(b.buf, byte(*n >> 24))
	b.buf = append(b.buf, byte(*n >> 32))
	b.buf = append(b.buf, byte(*n >> 40))
	b.buf = append(b.buf, byte(*n >> 48))
	b.buf = append(b.buf, byte(*n >> 56))
	b.pos += 4
	return
}

// Int8 reads or writes a int8
func (b *Stream) Int8(n *int8) {
	if b.err != nil {
		return
	}
	if b.reading {
		if b.pos+1 > len(b.buf) {
			b.err = io.EOF
			return
		}
		buf := b.buf[b.pos: b.pos+SizeInt8]
		b.pos += 1
		*n = int8(buf[0])
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.pos++
	return
}

// Int16 reads or writes a int16
func (b *Stream) Int16(n *int16) {
	if b.err != nil {
		return
	}
	if b.reading {
		buf := b.GetBytes(SizeInt16)
		if b.err != nil {
			return
		}
		var v int16
		v |= int16(buf[0])
		v |= int16(buf[1]) << 8
		*n = v
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.buf = append(b.buf, byte(*n >> 8))
	b.pos += 2
	return
}

// Int32 reads or writes a int32
func (b *Stream) Int32(n *int32) {
	if b.err != nil {
		return
	}
	if b.reading {
		buf := b.GetBytes(SizeInt32)
		if b.err != nil {
			return
		}
		var v int32
		v |= int32(buf[0])
		v |= int32(buf[1]) << 8
		v |= int32(buf[2]) << 16
		v |= int32(buf[3]) << 24
		*n = v
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.buf = append(b.buf, byte(*n >> 8))
	b.buf = append(b.buf, byte(*n >> 16))
	b.buf = append(b.buf, byte(*n >> 24))
	b.pos += 4
	return
}

// Int64 reads or writes a int64
func (b *Stream) Int64(n *int64) {
	if b.err != nil {
		return
	}
	if b.reading {
		buf := b.GetBytes(SizeInt64)
		if b.err != nil {
			return
		}
		var v int64
		v |= int64(buf[0])
		v |= int64(buf[1]) << 8
		v |= int64(buf[2]) << 16
		v |= int64(buf[3]) << 24
		v |= int64(buf[4]) << 32
		v |= int64(buf[5]) << 40
		v |= int64(buf[6]) << 48
		v |= int64(buf[7]) << 56
		*n = v
		return
	}
	b.buf = append(b.buf, byte(*n))
	b.buf = append(b.buf, byte(*n >> 8))
	b.buf = append(b.buf, byte(*n >> 16))
	b.buf = append(b.buf, byte(*n >> 24))
	b.buf = append(b.buf, byte(*n >> 32))
	b.buf = append(b.buf, byte(*n >> 40))
	b.buf = append(b.buf, byte(*n >> 48))
	b.buf = append(b.buf, byte(*n >> 56))
	b.pos += 4
	return
}

// Float32 reads or writes a float32
func (b *Stream) Float32(n *float32) {
	if b.err != nil {
		return
	}
	if b.reading {
		var v uint32
		b.Uint32(&v)
		*n = math.Float32frombits(v)
		return
	}
	var v = math.Float32bits(*n)
	b.Uint32(&v)
	return
}

// Float64 reads or writes a float64
func (b *Stream) Float64(n *float64) {
	if b.err != nil {
		return
	}
	if b.reading {
		var v uint64
		b.Uint64(&v)
		*n = math.Float64frombits(v)
		return
	}
	var v = math.Float64bits(*n)
	b.Uint64(&v)
	return
}

// GetByte decodes a little-endian byte
func (b *Stream) GetByte() (result byte) {
	if b.err != nil {
		return
	}
	var v uint8
	b.Uint8(&v)
	return v
}

// WriteByte encodes a little-endian uint8 into the buffer.
func (b *Stream) WriteByte(n byte) {
	b.buf = append(b.buf, uint8(n))
	b.pos++
}

// WriteBytes encodes a little-endian byte slice into the buffer
func (b *Stream) WriteBytes(src []byte) {
	for i := 0; i < len(src); i += 1 {
		b.WriteByte(uint8(src[i]))
	}
}

// WriteBytes encodes a little-endian byte slice into the buffer
func (b *Stream) WriteBytesN(src []byte, length int) {
	for i := 0; i < length; i += 1 {
		b.WriteByte(uint8(src[i]))
	}
}
