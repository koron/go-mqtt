package packet

import "testing"

func TestPingReq(t *testing.T) {
	data := []byte{0xc0, 0x00}
	p := PingReq{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.Type != TPingReq {
		t.Errorf("unexpected Type: %v", p.Type)
	}
	if p.Dup {
		t.Errorf("unexpected Dup: %v", p.Dup)
	}
	if p.QoS != QAtMostOnce {
		t.Errorf("unexpected QoS: %v", p.QoS)
	}
	if p.Retain {
		t.Errorf("unexpected Retain: %v", p.Retain)
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
	if p.Type != TPingResp {
		t.Errorf("unexpected Type: %v", p.Type)
	}
	if p.Dup {
		t.Errorf("unexpected Dup: %v", p.Dup)
	}
	if p.QoS != QAtMostOnce {
		t.Errorf("unexpected QoS: %v", p.QoS)
	}
	if p.Retain {
		t.Errorf("unexpected Retain: %v", p.Retain)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}
