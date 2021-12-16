package bug

import (
	"bufio"
)

// Document represents BSON Document data type.
type Document struct {
	m    map[string]any
	keys []string
}

// ReadFrom implements bsontype interface.
func (doc *Document) ReadFrom(r *bufio.Reader) error {
	r.Read(make([]byte, 5))

	doc.m = map[string]any{}
	doc.keys = []string{}
	return nil
}

// WriteTo implements bsontype interface.
func (doc Document) WriteTo(w *bufio.Writer) error {
	_, err := w.Write([]byte{0x05, 0x00, 0x00, 0x00, 0x00})
	return err
}

// MarshalBinary implements bsontype interface.
func (doc Document) MarshalBinary() ([]byte, error) {
	return []byte{0x05, 0x00, 0x00, 0x00, 0x00}, nil
}
