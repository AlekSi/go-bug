package bson

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
)

type bsontype interface {
	bsontype() // seal

	ReadFrom(*bufio.Reader) error
	WriteTo(*bufio.Writer) error
	encoding.BinaryMarshaler
	json.Unmarshaler
	json.Marshaler
}

//go-sumtype:decl bsontype

func checkConsumed(dec *json.Decoder, r *bytes.Reader) error {
	if dr := dec.Buffered().(*bytes.Reader); dr.Len() != 0 {
		b, _ := io.ReadAll(dr)
		return fmt.Errorf("%d bytes remains in the decoded: %s", dr.Len(), b)
	}

	if l := r.Len(); l != 0 {
		b, _ := io.ReadAll(r)
		return fmt.Errorf("%d bytes remains in the reader: %s", l, b)
	}

	return nil
}
