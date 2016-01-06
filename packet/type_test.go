package packet

import (
	"testing"
)

func TestTypeName(t *testing.T) {
	if s := TReserved.Name(); s != "RESERVED" {
		t.Errorf("TReserved.Name() returns %s", s)
	}
	if s := TConnect.Name(); s != "CONNECT" {
		t.Errorf("TConnect.Name() returns %s", s)
	}
	if s := TConnACK.Name(); s != "CONNACK" {
		t.Errorf("TConnACK.Name() returns %s", s)
	}
	if s := TPublish.Name(); s != "PUBLISH" {
		t.Errorf("TPublish.Name() returns %s", s)
	}
	if s := TPubACK.Name(); s != "PUBACK" {
		t.Errorf("TPubACK.Name() returns %s", s)
	}
	if s := TPubRec.Name(); s != "PUBREC" {
		t.Errorf("TPubRec.Name() returns %s", s)
	}
	if s := TPubRel.Name(); s != "PUBREL" {
		t.Errorf("TPubRel.Name() returns %s", s)
	}
	if s := TPubComp.Name(); s != "PUBCOMP" {
		t.Errorf("TPubComp.Name() returns %s", s)
	}
	if s := TSubscribe.Name(); s != "SUBSCRIBE" {
		t.Errorf("TSubscribe.Name() returns %s", s)
	}
	if s := TSubACK.Name(); s != "SUBACK" {
		t.Errorf("TSubACK.Name() returns %s", s)
	}
	if s := TUnsubscribe.Name(); s != "UNSUBSCRIBE" {
		t.Errorf("TUnsubscribe.Name() returns %s", s)
	}
	if s := TUnsubACK.Name(); s != "UNSUBACK" {
		t.Errorf("TUnsubACK.Name() returns %s", s)
	}
	if s := TPingReq.Name(); s != "PINGREQ" {
		t.Errorf("TPingReq.Name() returns %s", s)
	}
	if s := TPingResp.Name(); s != "PINGRESP" {
		t.Errorf("TPingResp.Name() returns %s", s)
	}
	if s := TDisconnect.Name(); s != "DISCONNECT" {
		t.Errorf("TDisconnect.Name() returns %s", s)
	}
	if s := TReserved2.Name(); s != "RESERVED2" {
		t.Errorf("TReserved2.Name() returns %s", s)
	}
}
