package main

import (
	"fmt"
	"net"
)

type HttpConnection struct {
	tcpConn net.Conn
	req     Request
	res     Response
}

func newHttpConnection(conn net.Conn, req Request) HttpConnection {
	return HttpConnection{tcpConn: conn, req: req, res: newResponse()}
}

func (hc *HttpConnection) Respond() error {
	_, err := hc.tcpConn.Write(hc.res.ToBytes())
	if err != nil {
		return fmt.Errorf("error writing response: %v", err.Error())
	}
	return nil
}
