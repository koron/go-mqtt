package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	"github.com/surgemq/message"
)

func writeMessage(w io.Writer, msg message.Message) error {
	b := make([]byte, msg.Len())
	_, err := msg.Encode(b)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return err
}

func writeConnackErrorMessage(w io.Writer, err error) error {
	cerr, ok := err.(message.ConnackCode);
	if !ok {
		return nil
	}
	resp := message.NewConnackMessage()
	resp.SetSessionPresent(false)
	resp.SetReturnCode(cerr)
	return writeMessage(w, resp)
}

func readRawMessage(r *bufio.Reader) ([]byte, error) {
	buf := &bytes.Buffer{}
	// read header: message type
	c, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	err = buf.WriteByte(c)
	if err != nil {
		return nil, err
	}
	// read length of payload.
	l, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	// FIXME: should be check l is too long.
	b := make([]byte, binary.MaxVarintLen32)
	n := binary.PutUvarint(b, l)
	buf.Write(b[0:n])
	// read whole payload
	_, err = io.CopyN(buf, r, int64(l))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func readConnectMessage(r *bufio.Reader) (*message.ConnectMessage, error) {
	raw, err := readRawMessage(r)
	if err != nil {
		return nil, err
	}
	msg := message.NewConnectMessage()
	_, err = msg.Decode(raw)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
