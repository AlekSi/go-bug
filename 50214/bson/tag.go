package bson

//go:generate ../../bin/stringer -linecomment -type tag

type tag byte

const (
	tagDocument = tag(0x03) // Document
	tagArray    = tag(0x04) // Array
)
