package client

import (
	"bytes"
	"errors"
)

var errBufTooSmall = errors.New("buffer is too small to fit the message")

const defaultScratchSize = 64 * 1024

//Client represents an instance of client connected to a set of Feuer servers.
type Simple struct {
	addrs []string

	buf     bytes.Buffer
	restBuf bytes.Buffer
}

// NewSimple creates a new client for the Feuer server.
func NewSimple(addrs []string) *Simple {
	return &Simple{
		addrs: addrs,
	}
}

// Send sends the messages to the Feuer servers.
func (s *Simple) Send(msgs []byte) error {
	_, err := s.buf.Write(msgs)
	return err
}

// Receive will either wait for new messages or return an
// error in case something goes wrong.
// The scratch buffer can be used to read the data.
func (s *Simple) Receive(scratch []byte) ([]byte, error) {
	if scratch == nil {
		scratch = make([]byte, defaultScratchSize)
	}

	n, err := s.buf.Read(scratch)
	if err != nil {
		return nil, err
	}

	truncated, rest, err := cutToLastMessage(scratch[0:n])
	if err != nil {
		return nil, err
	}

	if s.restBuf.Len() > 0 {
		truncated = append(s.restBuf.Bytes(), truncated...)
		s.restBuf.Reset()
	}

	if len(rest) != 0 {
		s.restBuf.Write(rest)
	}

	return truncated, nil
}

func cutToLastMessage(res []byte) (truncated []byte, rest []byte, err error) {
	n := len(res)

	if n == 0 {
		return res, nil, nil
	}

	if res[n-1] == '\n' {
		return res, nil, nil
	}

	lastPos := bytes.LastIndexByte(res, '\n')
	if lastPos < 0 {
		return nil, nil, errBufTooSmall
	}
	return res[0 : lastPos+1], res[lastPos+1:], nil
}
