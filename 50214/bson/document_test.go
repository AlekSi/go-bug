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
	"testing"

	"github.com/AlekSi/go-bug/50214/types"
	"github.com/AlekSi/go-bug/50214/util/testutil"
)

var (
	handshake1 = testCase{
		name: "handshake1",
		v: MustConvertDocument(types.MustMakeDocument(
			"ismaster", true,
			"client", types.MustMakeDocument(
				"driver", types.MustMakeDocument(
					"name", "nodejs",
					"version", "4.0.0-beta.6",
				),
				"os", types.MustMakeDocument(
					"type", "Darwin",
					"name", "darwin",
					"architecture", "x64",
					"version", "20.6.0",
				),
				"platform", "Node.js v14.17.3, LE (unified)|Node.js v14.17.3, LE (unified)",
				"application", types.MustMakeDocument(
					"name", "mongosh 1.0.1",
				),
			),
			"compression", types.Array{"none"},
			"loadBalanced", false,
		)),
		b: testutil.MustParseDumpFile("testdata", "handshake1.hex"),
		j: `{"$k":["ismaster","client","compression","loadBalanced"],"ismaster":true,` +
			`"client":{"$k":["driver","os","platform","application"],"driver":{"$k":["name","version"],` +
			`"name":"nodejs","version":"4.0.0-beta.6"},"os":{"$k":["type","name","architecture","version"],` +
			`"type":"Darwin","name":"darwin","architecture":"x64","version":"20.6.0"},` +
			`"platform":"Node.js v14.17.3, LE (unified)|Node.js v14.17.3, LE (unified)",` +
			`"application":{"$k":["name"],"name":"mongosh 1.0.1"}},"compression":["none"],"loadBalanced":false}`,
	}

	handshake2 = testCase{
		name: "handshake2",
		v: MustConvertDocument(types.MustMakeDocument(
			"ismaster", true,
			"client", types.MustMakeDocument(
				"driver", types.MustMakeDocument(
					"name", "nodejs",
					"version", "4.0.0-beta.6",
				),
				"os", types.MustMakeDocument(
					"type", "Darwin",
					"name", "darwin",
					"architecture", "x64",
					"version", "20.6.0",
				),
				"platform", "Node.js v14.17.3, LE (unified)|Node.js v14.17.3, LE (unified)",
				"application", types.MustMakeDocument(
					"name", "mongosh 1.0.1",
				),
			),
			"compression", types.Array{"none"},
			"loadBalanced", false,
		)),
		b: testutil.MustParseDumpFile("testdata", "handshake2.hex"),
		j: `{"$k":["ismaster","client","compression","loadBalanced"],"ismaster":true,` +
			`"client":{"$k":["driver","os","platform","application"],"driver":{"$k":["name","version"],` +
			`"name":"nodejs","version":"4.0.0-beta.6"},"os":{"$k":["type","name","architecture","version"],` +
			`"type":"Darwin","name":"darwin","architecture":"x64","version":"20.6.0"},` +
			`"platform":"Node.js v14.17.3, LE (unified)|Node.js v14.17.3, LE (unified)",` +
			`"application":{"$k":["name"],"name":"mongosh 1.0.1"}},"compression":["none"],"loadBalanced":false}`,
	}

	handshake3 = testCase{
		name: "handshake3",
		v: MustConvertDocument(types.MustMakeDocument(
			"lsid", types.MustMakeDocument(
				"id", types.Binary{
					Subtype: types.BinaryUUID,
					B:       []byte{0xa3, 0x19, 0xf2, 0xb4, 0xa1, 0x75, 0x40, 0xc7, 0xb8, 0xe7, 0xa3, 0xa3, 0x2e, 0xc2, 0x56, 0xbe},
				},
			),
			"$db", "admin",
		)),
		b: testutil.MustParseDumpFile("testdata", "handshake3.hex"),
		j: "{\"$k\":[\"lsid\",\"$db\"],\"lsid\":{\"$k\":[\"id\"],\"id\":{\"$b\":\"oxnytKF1QMe456OjLsJWvg==\",\"s\":4}},\"$db\":\"admin\"}",
	}

	documentTestCases = []testCase{handshake1, handshake2, handshake3}
)

func TestDocument(t *testing.T) {
	t.Parallel()

	t.Run("Binary", func(t *testing.T) {
		t.Parallel()
		testBinary(t, documentTestCases, func() bsontype { return new(Document) })
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		testJSON(t, documentTestCases, func() bsontype { return new(Document) })
	})
}

func FuzzDocumentBinary(f *testing.F) {
	fuzzBinary(f, documentTestCases, func() bsontype { return new(Document) })
}

func FuzzDocumentJSON(f *testing.F) {
	fuzzJSON(f, documentTestCases, func() bsontype { return new(Document) })
}

func BenchmarkDocument(b *testing.B) {
	benchmark(b, documentTestCases, func() bsontype { return new(Document) })
}
