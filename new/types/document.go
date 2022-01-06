package types

import (
	"fmt"
)

// Common interface with bson.Document.
type document interface {
	Map() map[string]any
	Keys() []string
}

// Document represents BSON document.
//
// Duplicate field names are not supported.
type Document struct {
	m    map[string]any
	keys []string
}

// ConvertDocument converts bson.Document to types.Document and validates it.
// It references the same data without copying it.
func ConvertDocument(d document) (Document, error) {
	if d == nil {
		panic("types.ConvertDocument: d is nil")
	}

	doc := Document{
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

// MakeDocument makes a new Document from given key/value pairs.
func MakeDocument(pairs ...any) (Document, error) {
	l := len(pairs)
	if l%2 != 0 {
		return Document{}, fmt.Errorf("types.MakeDocument: invalid number of arguments: %d", l)
	}

	doc := Document{
		m:    make(map[string]any, l/2),
		keys: make([]string, 0, l/2),
	}
	for i := 0; i < l; i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return Document{}, fmt.Errorf("types.MakeDocument: invalid key type: %T", pairs[i])
		}

		value := pairs[i+1]
		if err := doc.add(key, value); err != nil {
			return Document{}, fmt.Errorf("types.MakeDocument: %w", err)
		}
	}

	return doc, nil
}

// MustMakeDocument is a MakeDocument that panics in case of error.
func MustMakeDocument(pairs ...any) Document {
	doc, err := MakeDocument(pairs...)
	if err != nil {
		panic(err)
	}
	return doc
}

// Map returns a shallow copy of the document as a map. Do not modify it.
func (d Document) Map() map[string]any {
	return d.m
}

// Keys returns a shallow copy of the document's keys. Do not modify it.
func (d Document) Keys() []string {
	return d.keys
}

func (d *Document) add(key string, value any) error {
	if _, ok := d.m[key]; ok {
		return fmt.Errorf("types.Document.add: key already present: %q", key)
	}

	d.keys = append(d.keys, key)
	d.m[key] = value

	return nil
}

// check interfaces
var (
	_ document = Document{}
	_ document = (*Document)(nil)
)
