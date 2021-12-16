package bson

import (
	"bufio"
	"encoding"
)

type bsontype interface {
	bsontype() // seal

	ReadFrom(*bufio.Reader) error
	WriteTo(*bufio.Writer) error
	encoding.BinaryMarshaler
}

//go-sumtype:decl bsontype
