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

//go:generate ../../bin/stringer -linecomment -type tag

type tag byte

const (
	tagDouble   = tag(0x01) // Double
	tagString   = tag(0x02) // String
	tagDocument = tag(0x03) // Document
	tagArray    = tag(0x04) // Array
	tagBinary   = tag(0x05) // Binary
	tagBool     = tag(0x08) // Bool
)
