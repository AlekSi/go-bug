package bson

import (
	"bufio"
)

const (
	MaxDocumentLen = 16777216

	minDocumentLen = 5
)

// Common interface with types.Document.
type document interface {
	Map() map[string]any
	Keys() []string
}

// Document represents BSON Document data type.
type Document struct {
	m    map[string]any
	keys []string
}

// convertDocument converts types.Document to bson.Document and validates it.
// It references the same data without copying it.
func mustConvertDocument(d document) *Document {
	doc := &Document{
		m:    d.Map(),
		keys: d.Keys(),
	}

	if doc.m == nil {
		doc.m = map[string]any{}
	}
	if doc.keys == nil {
		doc.keys = []string{}
	}

	return doc
}

func (doc *Document) bsontype() {}

// Returns the map of key values associated with the Document.
func (doc *Document) Map() map[string]any {
	return doc.m
}

// Keys returns the keys associated with the document.
func (doc *Document) Keys() []string {
	return doc.keys
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

// check interfaces
var (
	_ bsontype = (*Document)(nil)
	_ document = (*Document)(nil)
)
