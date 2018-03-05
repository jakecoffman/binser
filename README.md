# binser

Serializer helper inspired by https://gafferongames.com/post/serialization_strategies/

## about

Tired of having to write a separate binary marshaler and unmarshaler for each of your packet types?
Or maybe you are just tired of them getting out of sync, causing minutes of tracking down a new issue
only to find out it was this again.

Here's an example of how to use it:

```Go
func (l Location) MarshalBinary() ([]byte, error) {
	return l.Serialize(nil)
}

func (l *Location) UnmarshalBinary(b []byte) error {
	_, err := l.Serialize(b)
	return err
}

func (l *Location) Serialize(b []byte) ([]byte, error) {
	stream := binser.NewStream(b)
	var m uint8 = PacketLocation
	stream.Uint8(&m)
	stream.Uint16((*uint16)(&l.ID))
	stream.Uint64(&l.Sequence)
	stream.Float32(&l.X)
	stream.Float32(&l.Y)
	stream.Float32(&l.Vx)
	stream.Float32(&l.Vy)
	stream.Float32(&l.Angle)
	stream.Float32(&l.AngularVelocity)
	stream.Float32(&l.Turret)
	return stream.Bytes()
}
```

## documentation

[godocs](https://godoc.org/github.com/jakecoffman/binser) 

## originally forked from

[Isaac Dawson](https://github.com/wirepair/binserializer)
