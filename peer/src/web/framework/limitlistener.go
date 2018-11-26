package framework

import (
	"net"
	"sync"
	"time"
)

/*
结合echo.go中的tcpKeepAliveListener及netutils中listen.go的
LimitListener进行如下设计，能够限制http连接数
*/

type limitListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release     func()
}

func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)
	return err
}

type tcpKeepAliveListener struct {
	*net.TCPListener
	sem chan struct{}
}

func (ln *tcpKeepAliveListener) acquire() { ln.sem <- struct{}{} }
func (ln *tcpKeepAliveListener) release() { <-ln.sem }

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	ln.acquire()
	tc, err := ln.AcceptTCP()
	if err != nil {
		ln.release()
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Second)
	return &limitListenerConn{Conn: tc, release: ln.release}, nil
}

func newLimitListener(address string, n int) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener), make(chan struct{}, n)}, nil
}
