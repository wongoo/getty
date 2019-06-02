package tcp

import (
	"fmt"
	"github.com/dubbogo/getty"
	"github.com/dubbogo/getty/examples/hello"

	"net"
	"time"
)

func NewHelloServerSession(session getty.Session) (err error) {
	return InitialSession(session, hello.NewHelloPackageHandler(), hello.NewHelloMessageListener())
}

func NewHelloClientSession(session getty.Session) (err error) {
	listener := hello.NewHelloMessageListener()
	listener.SessionOnOpen = func(session getty.Session) {
		hello.Sessions = append(hello.Sessions, session)
	}
	return InitialSession(session, hello.NewHelloPackageHandler(), listener)
}

func InitialSession(session getty.Session, pkgHandler getty.ReadWriter, listener getty.EventListener) (err error) {
	session.SetCompressType(getty.CompressZip)

	tcpConn, ok := session.Conn().(*net.TCPConn)
	if !ok {
		panic(fmt.Sprintf("newSession: %s, session.conn{%#v} is not tcp connection", session.Stat(), session.Conn()))
	}

	if err = tcpConn.SetNoDelay(true); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlive(true); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlivePeriod(10 * time.Second); err != nil {
		return err
	}
	if err = tcpConn.SetReadBuffer(262144); err != nil {
		return err
	}
	if err = tcpConn.SetWriteBuffer(524288); err != nil {
		return err
	}

	session.SetName("hello")
	session.SetMaxMsgLen(128)
	session.SetRQLen(1024)
	session.SetWQLen(512)
	session.SetReadTimeout(time.Second)
	session.SetWriteTimeout(5 * time.Second)
	session.SetCronPeriod(int(hello.CronPeriod.Nanoseconds() / 1e6))
	session.SetWaitTime(time.Second)

	session.SetPkgHandler(pkgHandler)
	session.SetEventListener(listener)
	return nil
}
