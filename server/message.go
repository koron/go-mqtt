package server

import (
	"io"

	"github.com/koron/go-mqtt/packet"
)

func writeMessage(w io.Writer, p packet.Packet) (int, error) {
	b, err := p.Encode()
	if err != nil {
		return 0, err
	}
	return w.Write(b)
}

func writeConnackErrorMessage(w io.Writer, err error) (int, error) {
	// TODO: rewrite writeConnackErrorMessage
	return 0, nil
}

func readConnectMessage(r packet.Reader) (*packet.Connect, error) {
	b, err := packet.Split(r)
	if err != nil {
		return nil, err
	}
	p := packet.Connect{}
	if err := p.Decode(b); err != nil {
		return nil, err
	}
	return &p, nil
}
