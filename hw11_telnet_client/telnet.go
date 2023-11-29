package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrConnClose = errors.New("connection was closed")
	ErrEOF       = errors.New("EOF")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address     string
	timeout     time.Duration
	in          io.ReadCloser
	out         io.Writer
	conn        net.Conn
	connScanner *bufio.Scanner
	inScanner   *bufio.Scanner
}

func NewTelnetClient(
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	t.connScanner = bufio.NewScanner(t.conn)
	t.inScanner = bufio.NewScanner(t.in)

	return nil
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return nil
	}
	if !t.connScanner.Scan() {
		return ErrConnClose
	}
	_, err := t.out.Write([]byte(t.connScanner.Text() + "\n"))
	return err
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return nil
	}
	if !t.inScanner.Scan() {
		return ErrEOF
	}
	_, err := t.conn.Write([]byte(t.inScanner.Text() + "\n"))
	return err
}

func (t *telnetClient) Close() error {
	var err error
	if t.conn != nil {
		err = t.conn.Close()
	}
	if t.in != nil {
		err = t.in.Close()
	}
	return err
}
