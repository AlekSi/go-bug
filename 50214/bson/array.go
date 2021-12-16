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
	"encoding/json"
	"fmt"
	"strconv"
)

// Array represents BSON Array data type.
type Array []any

func (arr *Array) bsontype() {}

// ReadFrom implements bsontype interface.
func (arr *Array) ReadFrom(r *bufio.Reader) error {
	var doc Document
	if err := doc.ReadFrom(r); err != nil {
		return err
	}

	s := make([]any, len(doc.m))

	for i := 0; i < len(doc.m); i++ {
		if k := doc.keys[i]; k != strconv.Itoa(i) {
			return fmt.Errorf("key %d is %q", i, k)
		}

		v, ok := doc.m[strconv.Itoa(i)]
		if !ok {
			return fmt.Errorf("no element %d in array of length %d", i, len(doc.m))
		}
		s[i] = v
	}

	*arr = s

	return nil
}

// WriteTo implements bsontype interface.
func (arr Array) WriteTo(w *bufio.Writer) error {
	v, err := arr.MarshalBinary()
	if err != nil {
		return err
	}

	if _, err = w.Write(v); err != nil {
		return err
	}

	return nil
}

// MarshalBinary implements bsontype interface.
func (arr Array) MarshalBinary() ([]byte, error) {
	m := make(map[string]any, len(arr))
	keys := make([]string, len(arr))
	for i := 0; i < len(keys); i++ {
		key := strconv.Itoa(i)
		m[key] = arr[i]
		keys[i] = key
	}

	doc := Document{
		m:    m,
		keys: keys,
	}
	b, err := doc.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UnmarshalJSON implements bsontype interface.
func (arr *Array) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)

	var rawMessages []json.RawMessage
	if err := dec.Decode(&rawMessages); err != nil {
		return err
	}
	if err := checkConsumed(dec, r); err != nil {
		return err
	}

	*arr = make(Array, len(rawMessages))
	for i, el := range rawMessages {
		v, err := unmarshalJSONValue(el)
		if err != nil {
			return err
		}

		(*arr)[i] = v
	}

	return nil
}

// MarshalJSON implements bsontype interface.
func (arr Array) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')

	for i, el := range arr {
		if i != 0 {
			buf.WriteByte(',')
		}

		b, err := marshalJSONValue(el)
		if err != nil {
			return nil, err
		}

		buf.Write(b)
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// check interfaces
var (
	_ bsontype = (*Array)(nil)
)
