package packet

import (
	"testing"
)

func TestTypeName(t *testing.T) {
	for _, tc := range []struct {
		typ  Type
		name string
	}{
		{TReserved, "RESERVED"},
		{TConnect, "CONNECT"},
		{TConnACK, "CONNACK"},
		{TPublish, "PUBLISH"},
		{TPubACK, "PUBACK"},
		{TPubRec, "PUBREC"},
		{TPubRel, "PUBREL"},
		{TPubComp, "PUBCOMP"},
		{TSubscribe, "SUBSCRIBE"},
		{TSubACK, "SUBACK"},
		{TUnsubscribe, "UNSUBSCRIBE"},
		{TUnsubACK, "UNSUBACK"},
		{TPingReq, "PINGREQ"},
		{TPingResp, "PINGRESP"},
		{TDisconnect, "DISCONNECT"},
		{TReserved2, "RESERVED2"},
	} {
		if s := tc.typ.Name(); s != tc.name {
			t.Errorf("mismatch: expected=%s actual=%s", tc.name, s)
		}
	}
}
