package netUtils

import (
	"net"
	"reflect"
)

func getFdFromTCPConn(tcpConn *net.TCPConn) int {
	conn := reflect.ValueOf(*tcpConn)
	valueConn := conn.FieldByName("conn")
	fd := valueConn.FieldByName("fd")
	ptr := reflect.Indirect(fd)
	connFd := ptr.FieldByName("sysfd")
	return int(connFd.Int())
}
