package packet

import "testing"

func compareTopics(t *testing.T, actual, expected []Topic) {
	la, lb := len(actual), len(expected)
	n := min(la, lb)
	for i := 0; i < n; i++ {
		if actual[i].Filter != expected[i].Filter {
			t.Errorf("Topic.Filter#%d isn't match: expected=%s actual=%s",
				i, expected[i].Filter, actual[i].Filter)
			return
		}
		if actual[i].RequestedQoS != expected[i].RequestedQoS {
			t.Errorf("Topic.RequestedQoS#%d isn't match: expected=%v actual=%v",
				i, expected[i].RequestedQoS, actual[i].RequestedQoS)
			return
		}
	}
	if la != lb {
		if la > lb {
			t.Errorf("len(actual)=%d > len(expected)=%d actual[%d]=%02x",
				la, lb, lb, actual[lb])
		} else {
			t.Errorf("len(actual)=%d < len(expected)=%d expected[%d]=%02x",
				la, lb, la, expected[la])
		}
	}
}

func compareSubscribeResults(t *testing.T, expected, actual []SubscribeResult) {
	la, lb := len(actual), len(expected)
	n := min(la, lb)
	for i := 0; i < n; i++ {
		if actual[i] != expected[i] {
			t.Errorf("SubscribeResult#%d isn't match: expected=%v actual=%v",
				i, expected[i], actual[i])
			return
		}
	}
	if la != lb {
		if la > lb {
			t.Errorf("len(actual)=%d > len(expected)=%d actual[%d]=%v",
				la, lb, lb, actual[lb])
		} else {
			t.Errorf("len(actual)=%d < len(expected)=%d expected[%d]=%v",
				la, lb, la, expected[la])
		}
	}
}

func TestSubscribe(t *testing.T) {
	data := []byte{
		0x82,
		36,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		'g', 'o', '-', 'm', 'q', 't', 't',
		0, // QoS
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}
	p := Subscribe{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %v", p.PacketID)
	}
	compareTopics(t, p.Topics, []Topic{
		{Filter: "go-mqtt", RequestedQoS: QAtMostOnce},
		{Filter: "/a/b/#/c", RequestedQoS: QAtLeastOnce},
		{Filter: "/a/b/#/cdd", RequestedQoS: QExactlyOnce},
	})

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestSubACK(t *testing.T) {
	data := []byte{
		0x90,
		6,
		0,    // packet ID MSB (0)
		7,    // packet ID LSB (7)
		0,    // return code 1
		1,    // return code 2
		2,    // return code 3
		0x80, // return code 4
	}
	p := SubACK{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}
	compareSubscribeResults(t, p.Results, []SubscribeResult{
		SubscribeAtMostOnce,
		SubscribeAtLeastOnce,
		SubscribeExactOnce,
		SubscribeFailure,
	})

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestUnsubscribe(t *testing.T) {
	data := []byte{
		0xa2,
		33,
		0, // packet ID MSB (0)
		7, // packet ID LSB (7)
		0, // topic name MSB (0)
		7, // topic name LSB (7)
		'g', 'o', '-', 'm', 'q', 't', 't',
		0, // topic name MSB (0)
		8, // topic name LSB (8)
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		0,  // topic name MSB (0)
		10, // topic name LSB (10)
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
	}
	p := Unsubscribe{}
	err := p.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if p.PacketID != 7 {
		t.Errorf("unexpected PacketID: %d", p.PacketID)
	}
	compareStrings(t, p.Topics, []string{
		"go-mqtt", "/a/b/#/c", "/a/b/#/cdd",
	})

	// encode test.
	b, err := p.Encode()
	if err != nil {
		t.Fatal(err)
	}
	compareBytes(t, b, data)
}

func TestUnsubACK(t *testing.T) {
	data := []byte{
		0xb0,
		0x02,
		0x00,
		0x07,
	}
	p := UnsubACK{}
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
