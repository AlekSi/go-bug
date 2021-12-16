// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bson

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
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

	// for validation
	if _, err := types.ConvertDocument(doc); err != nil {
		return nil, fmt.Errorf("bson.ConvertDocument: %w", err)
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
		case tagString:
			var v String
			if err := v.ReadFrom(bufr); err != nil {
				return fmt.Errorf("bson.Document.ReadFrom (String): %w", err)
			}
			doc.m[string(ename)] = string(v)

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
		case string:
			bufw.WriteByte(byte(tagString))
			if err := ename.WriteTo(bufw); err != nil {
				return nil, err
			}
			if err := String(elV).WriteTo(bufw); err != nil {
				return nil, err
			}

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

func unmarshalJSONValue(data []byte) (any, error) {
	var v any
	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	err := dec.Decode(&v)
	if err != nil {
		return nil, err
	}
	if err := checkConsumed(dec, r); err != nil {
		return nil, err
	}

	var res any
	switch v := v.(type) {
	case map[string]any:
		switch {
		case v["$k"] != nil:
			var o Document
			err = o.UnmarshalJSON(data)
			if err == nil {
				res, err = types.ConvertDocument(&o)
			}
		default:
			err = fmt.Errorf("unmarshalJSONValue: unhandled map %v", v)
		}
	case string:
		res = v
	case []any:
		var o Array
		err = o.UnmarshalJSON(data)
		res = types.Array(o)
	default:
		err = fmt.Errorf("unmarshalJSONValue: unhandled element %[1]T (%[1]v)", v)
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UnmarshalJSON implements bsontype interface.
func (doc *Document) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)

	var rawMessages map[string]json.RawMessage
	if err := dec.Decode(&rawMessages); err != nil {
		return err
	}
	if err := checkConsumed(dec, r); err != nil {
		return err
	}

	b, ok := rawMessages["$k"]
	if !ok {
		return fmt.Errorf("bson.Document.UnmarshalJSON: missing $k")
	}

	var keys []string
	if err := json.Unmarshal(b, &keys); err != nil {
		return err
	}
	if len(keys)+1 != len(rawMessages) {
		return fmt.Errorf("bson.Document.UnmarshalJSON: %d elements in $k, %d in total", len(keys), len(rawMessages))
	}

	doc.keys = keys
	doc.m = make(map[string]any, len(keys))

	for _, key := range keys {
		b, ok = rawMessages[key]
		if !ok {
			return fmt.Errorf("bson.Document.UnmarshalJSON: missing key %q", key)
		}
		v, err := unmarshalJSONValue(b)
		if err != nil {
			return err
		}
		doc.m[key] = v
	}

	if _, err := types.ConvertDocument(doc); err != nil {
		return fmt.Errorf("bson.Document.UnmarshalJSON: %w", err)
	}

	return nil
}

func marshalJSONValue(v any) ([]byte, error) {
	var o json.Marshaler
	var err error
	switch v := v.(type) {
	case string:
		o = String(v)
	case types.Document:
		o, err = ConvertDocument(v)
	case types.Array:
		o = Array(v)
	case nil:
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("marshalJSONValue: unhandled type %T", v)
	}

	if err != nil {
		return nil, err
	}

	b, err := o.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalJSON implements bsontype interface.
func (doc Document) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`{"$k":`)
	b, err := json.Marshal(doc.keys)
	if err != nil {
		return nil, err
	}
	buf.Write(b)

	for _, key := range doc.keys {
		buf.WriteByte(',')

		if b, err = json.Marshal(key); err != nil {
			return nil, err
		}
		buf.Write(b)
		buf.WriteByte(':')

		value := doc.m[key]
		b, err := marshalJSONValue(value)
		if err != nil {
			return nil, fmt.Errorf("bson.Document.MarshalJSON: %w", err)
		}

		buf.Write(b)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// check interfaces
var (
	_ bsontype = (*Document)(nil)
	_ document = (*Document)(nil)
)
