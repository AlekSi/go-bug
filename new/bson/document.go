package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/AlekSi/go-bug/new/types"
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

// ConvertDocument converts types.Document to bson.Document and validates it.
// It references the same data without copying it.
func ConvertDocument(d document) (*Document, error) {
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

	return doc, nil
}

// MustConvertDocument is a ConvertDocument that panics in case of error.
func MustConvertDocument(d document) *Document {
	doc, err := ConvertDocument(d)
	if err != nil {
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

		var ename CString
		if err := ename.ReadFrom(bufr); err != nil {
			return fmt.Errorf("bson.Document.ReadFrom (ename.ReadFrom): %w", err)
		}

		doc.keys = append(doc.keys, string(ename))

		switch tag(t) {
		case tagDocument:
			// TODO check maximum nesting

			var v Document
			if err := v.ReadFrom(bufr); err != nil {
				return fmt.Errorf("bson.Document.ReadFrom (embedded document): %w", err)
			}
			doc.m[string(ename)], err = types.ConvertDocument(&v)
			if err != nil {
				return fmt.Errorf("bson.Document.ReadFrom (embedded document): %w", err)
			}

		case tagArray:
			// TODO check maximum nesting

			var v Array
			if err := v.ReadFrom(bufr); err != nil {
				return fmt.Errorf("bson.Document.ReadFrom (Array): %w", err)
			}
			doc.m[string(ename)] = types.Array(v)

		default:
			return fmt.Errorf("bson.Document.ReadFrom: unhandled element type %#02x", t)
		}
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
		return fmt.Errorf("bson.Document.WriteTo: %w", err)
	}

	_, err = w.Write(v)
	if err != nil {
		return fmt.Errorf("bson.Document.WriteTo: %w", err)
	}

	return nil
}

// MarshalBinary implements bsontype interface.
func (doc Document) MarshalBinary() ([]byte, error) {
	var elist bytes.Buffer
	bufw := bufio.NewWriter(&elist)

	for _, elK := range doc.keys {
		ename := CString(elK)
		elV, ok := doc.m[elK]
		if !ok {
			panic(fmt.Sprintf("%q not found in map", elK))
		}

		switch elV := elV.(type) {
		case types.Document:
			bufw.WriteByte(byte(tagDocument))
			if err := ename.WriteTo(bufw); err != nil {
				return nil, err
			}
			doc, err := ConvertDocument(elV)
			if err != nil {
				return nil, err
			}
			if err := doc.WriteTo(bufw); err != nil {
				return nil, err
			}

		case types.Array:
			bufw.WriteByte(byte(tagArray))
			if err := ename.WriteTo(bufw); err != nil {
				return nil, err
			}
			if err := Array(elV).WriteTo(bufw); err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("bson.Document.MarshalBinary: unhandled element type %T", elV)
		}
	}

	if err := bufw.Flush(); err != nil {
		return nil, err
	}

	var res bytes.Buffer
	l := int32(elist.Len() + 5)
	binary.Write(&res, binary.LittleEndian, l)
	if _, err := elist.WriteTo(&res); err != nil {
		panic(err)
	}
	res.WriteByte(0)
	if int32(res.Len()) != l {
		panic(fmt.Sprintf("got %d, expected %d", res.Len(), l))
	}
	return res.Bytes(), nil
}

// check interfaces
var (
	_ bsontype = (*Document)(nil)
	_ document = (*Document)(nil)
)
