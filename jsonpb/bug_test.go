package jsonpb

import (
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
)

func TestDuration(t *testing.T) {
	d := -time.Nanosecond
	dp := ptypes.DurationProto(d)
	m := &jsonpb.Marshaler{}
	s, err := m.MarshalToString(dp)
	if err != nil {
		t.Fatal(err)
	}
	if s != `"-0.00000001s"` {
		t.Errorf("Unexpected result: %s", s)
	}
}
