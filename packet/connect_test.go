package packet

import "testing"

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
	if p.ClientID != "go-mqtt" {
		t.Errorf("unexpected ClientID: %s", p.ClientID)
	}
	if p.Version != 4 {
		t.Errorf("unexpected Version: %d", p.Version)
	}
	if p.Username == nil || *p.Username != "username" {
		t.Errorf("unexpected Username: %v", p.Username)
	}
	if p.Password == nil || *p.Password != "verysecret" {
		t.Errorf("unexpected Password: %v", p.Password)
	}
	if !p.CleanSession {
		t.Errorf("unexpected CleanSession: %v", p.CleanSession)
	}
	if p.KeepAlive != 10 {
		t.Errorf("unexpected KeepAlive: %d", p.KeepAlive)
	}
	if !p.WillFlag {
		t.Errorf("unexpected WillFlag: %v", p.WillFlag)
	}
	if p.WillQoS != QAtLeastOnce {
		t.Errorf("unexpected WillQoS: %v", p.WillQoS)
	}
	if p.WillRetain {
		t.Errorf("unexpected WillRetain: %v", p.WillRetain)
	}
	if p.WillTopic != "will" {
		t.Errorf("unexpected WillTopic: %s", p.WillTopic)
	}
	if p.WillMessage != "send me home" {
		t.Errorf("unexpected WillMessage: %s", p.WillMessage)
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
