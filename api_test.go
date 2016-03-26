package openzwave

import "testing"

func Test_RoundtripMarshaling(t *testing.T) {
	a := &api{}
	a.init(a)
	c := a.C()
	aa := unmarshal(c).Go().(*api)
	if a != aa {
		t.Fatalf("failed to round trip")
	}
}
