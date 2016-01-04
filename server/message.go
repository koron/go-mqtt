package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"

	"github.com/surgemq/message"
)

func writeMessage(w io.Writer, msg message.Message) (int, error) {
	b := make([]byte, msg.Len())
	_, err := msg.Encode(b)
	if err != nil {
		return 0, err
	}
	return w.Write(b)
}

func writeConnackErrorMessage(w io.Writer, err error) (int, error) {
	cerr, ok := err.(message.ConnackCode)
	if !ok {
		return 0, nil
	}
	resp := message.NewConnackMessage()
	resp.SetSessionPresent(false)
	resp.SetReturnCode(cerr)
	return writeMessage(w, resp)
}

func readMessage(r *bufio.Reader) (message.Message, error) {
	b, err := readRawMessage(r)
	if err != nil {
		return nil, err
	}
	t := message.MessageType(b[0] >> 4)
	msg, err := t.New()
	if err != nil {
		return nil, err
	}
	_, err = msg.Decode(b)
	if err != nil {
		return nil, err
	}
	return msg, nil
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
