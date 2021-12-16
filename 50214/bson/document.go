package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/AlekSi/go-bug/50214/types"
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

	// for validation
	if _, err := types.ConvertDocument(doc); err != nil {
		panic(err)
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
	var l int32
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		return fmt.Errorf("bson.Document.ReadFrom (binary.Read): %w", err)
	}
	if l < minDocumentLen || l > MaxDocumentLen {
		return fmt.Errorf("bson.Document.ReadFrom: invalid length %d", l)
	}

	// make buffer
	b := make([]byte, l)

	binary.LittleEndian.PutUint32(b, uint32(l))

	// read e_list and terminating zero
	n, err := io.ReadFull(r, b[4:])
	if err != nil {
		return fmt.Errorf("bson.Document.ReadFrom (io.ReadFull, expected %d, read %d): %w", len(b), n, err)
	}

	bufr := bufio.NewReader(bytes.NewReader(b[4:]))
	doc.m = map[string]any{}
	doc.keys = make([]string, 0, 2)

	for {
		t, err := bufr.ReadByte()
		if err != nil {
			return fmt.Errorf("bson.Document.ReadFrom (ReadByte): %w", err)
		}

		if t == 0 {
			// documented ended
			if _, err := bufr.Peek(1); err != io.EOF {
				return fmt.Errorf("unexpected end of the document: %w", err)
			}
			break
		}

		panic("not reached")
	}

	if _, err := types.ConvertDocument(doc); err != nil {
		return fmt.Errorf("bson.Document.ReadFrom: %w", err)
	}

	return nil
}

// WriteTo implements bsontype interface.
func (doc Document) WriteTo(w *bufio.Writer) error {
	v, err := doc.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = w.Write(v)
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
