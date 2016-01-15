package packet

import "testing"

func TestPingReq(t *testing.T) {
	data := []byte{0xc0, 0x00}
	p := PingReq{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestPingResp(t *testing.T) {
	data := []byte{0xd0, 0x00}
	p := PingResp{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}
