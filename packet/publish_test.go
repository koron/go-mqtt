package packet

import "testing"

func TestPublish(t *testing.T) {
	data := []byte{
		0x3d, 0x17,
		0x00, 0x07,
		'g', 'o', '-', 'm', 'q', 't', 't',
		0x00, 0x07,
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
	}
	p := Publish{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.TopicName != "go-mqtt" {
		t.Errorf("unexpected TopicName: %q", p.TopicName)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}
	if string(p.Payload) != "send me home" {
		t.Errorf("unexpected Payload: %v", p.Payload)
	}
	if !p.Dup {
		t.Errorf("unexpected Dup: %v", p.Dup)
	}
	if p.QoS != QExactlyOnce {
		t.Errorf("unexpected QoS: %v", p.QoS)
	}
	if !p.Retain {
		t.Errorf("unexpected Retain: %v", p.Retain)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestPubACK(t *testing.T) {
	data := []byte{0x40, 0x02, 0x00, 0x07}
	p := PubACK{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestPubRec(t *testing.T) {
	data := []byte{0x50, 0x02, 0x00, 0x07}
	p := PubRec{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestPubRel(t *testing.T) {
	data := []byte{0x62, 0x02, 0x00, 0x07}
	p := PubRel{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.QoS != QAtLeastOnce {
		t.Errorf("unexpected QoS: %v", p.QoS)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestPubComp(t *testing.T) {
	data := []byte{0x70, 0x02, 0x00, 0x07}
	p := PubComp{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}
