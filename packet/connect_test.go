package packet

import (
	"reflect"
	"testing"
)

func TestConnect1(t *testing.T) {
	data := []byte{
		0x10,
		61,
		0, // Length MSB (0)
		4, // Length LSB (4)
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // connect flags 11001110, will QoS = 01
		0,   // Keep Alive MSB (0)
		10,  // Keep Alive LSB (10)
		0,   // Client ID MSB (0)
		7,   // Client ID LSB (7)
		'g', 'o', '-', 'm', 'q', 't', 't',
		0, // Will Topic MSB (0)
		4, // Will Topic LSB (4)
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB (0)
		12, // Will Message LSB (12)
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB (0)
		8, // Username ID LSB (7)
		'u', 's', 'e', 'r', 'n', 'a', 'm', 'e',
		0,  // Password ID MSB (0)
		10, // Password ID LSB (10)
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}
	p := Connect{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(Connect{
		ClientID:     "go-mqtt",
		Version:      4,
		Username:     str2ptr("username"),
		Password:     str2ptr("verysecret"),
		CleanSession: true,
		KeepAlive:    10,
		WillFlag:     true,
		WillQoS:      QAtLeastOnce,
		WillRetain:   false,
		WillTopic:    "will",
		WillMessage:  "send me home",
	}, p) {
		t.Fatalf("mismatch Connect:\n actual=%+v", p)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestConnACK1(t *testing.T) {
	data := []byte{0x20, 0x02, 0x01, 0x05}
	p := ConnACK{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if !p.SessionPresent {
		t.Errorf("unexpected SessionPresent: %v", p.SessionPresent)
	}
	if p.ReturnCode != ConnectNotAuthorized {
		t.Errorf("unexpected ReturnCode: %v", p.ReturnCode)
	}

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestDisconnect1(t *testing.T) {
	data := []byte{0xe0, 0x00}
	p := Disconnect{}
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
