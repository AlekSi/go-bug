package bson

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
)

// CString represents BSON CString data type.
type CString string

func (cstr *CString) bsontype() {}

// ReadFrom implements bsontype interface.
func (cstr *CString) ReadFrom(r *bufio.Reader) error {
	b, err := r.ReadBytes(0)
	if err != nil {
		return fmt.Errorf("bson.CString.ReadFrom: %w", err)
	}

	*cstr = CString(b[:len(b)-1])
	return nil
}

// WriteTo implements bsontype interface.
func (cstr CString) WriteTo(w *bufio.Writer) error {
	v, err := cstr.MarshalBinary()
	if err != nil {
		return fmt.Errorf("bson.CString.WriteTo: %w", err)
	}

	_, err = w.Write(v)
	if err != nil {
		return fmt.Errorf("bson.CString.WriteTo: %w", err)
	}

	return nil
}

// MarshalBinary implements bsontype interface.
func (cstr CString) MarshalBinary() ([]byte, error) {
	b := make([]byte, len(cstr)+1)
	copy(b, cstr)
	return b, nil
}

type cstringJSON struct {
	CString string `json:"$c"`
}

// UnmarshalJSON implements bsontype interface.
func (cstr *CString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		panic("null data")
	}

	r := bytes.NewReader(data)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var o cstringJSON
	if err := dec.Decode(&o); err != nil {
		return err
	}
	if err := checkConsumed(dec, r); err != nil {
		return fmt.Errorf("bson.CString.UnmarshalJSON: %s", err)
	}

	*cstr = CString(o.CString)
	return nil
}

// MarshalJSON implements bsontype interface.
func (cstr CString) MarshalJSON() ([]byte, error) {
	return json.Marshal(cstringJSON{
		CString: string(cstr),
	})
}

// check interfaces
var (
	_ bsontype = (*CString)(nil)
)
