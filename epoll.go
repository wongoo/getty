package getty

import (
	"net"
	"reflect"
)

var (
	epollBatchSize = 100
)

type epoller interface {
	add(ss *session) error
	remove(ss *session) error
	wait() ([]*session, error)
	start()
}

//SetEpollBatchSize set epoll batch size
func SetEpollBatchSize(size int) {
	epollBatchSize = size
}

// Extract file descriptor associated with the connection
func socketFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
