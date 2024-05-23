package main

import (
	"flag"
	"fmt"
	"net"
)

type Config struct {
	DirectoryFlag string
}

func newConfig() Config {
	directory := flag.String("directory", "", "Directory path")
	flag.Parse()
	return Config{DirectoryFlag: *directory}
}

type HttpConnection struct {
	tcpConn net.Conn
	req     Request
	res     Response
	config  Config
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
