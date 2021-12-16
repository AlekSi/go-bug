package bson

//go:generate ../../bin/stringer -linecomment -type tag

type tag byte

const (
	tagDocument = tag(0x03) // Document
)
