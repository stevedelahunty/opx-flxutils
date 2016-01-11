package ipcutils

import (
	//"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"models"
)

type IPCClientBase struct {
	Address               string
	TTransport            thrift.TTransport
	PtrProtocolFactory   *thrift.TBinaryProtocolFactory
	IsConnected           bool
}

func (clnt *IPCClientBase) IsConnectedToServer() bool {
	return clnt.IsConnected
}

func (clnt *IPCClientBase) GetBulkObject(obj models.ConfigObj, currMarker int64, count int64) (err error,
	objCount int64,
	nextMarker int64,
	more bool,
	objs []models.ConfigObj) {
	//logger.Println("### Get Bulk request called with", currMarker, count)
	return nil, 0, 0, false, make([]models.ConfigObj, 0)
}

//
// This method gets Thrift related IPC handles.
//
func CreateIPCHandles(address string) (thrift.TTransport, *thrift.TBinaryProtocolFactory) {
	var transportFactory thrift.TTransportFactory
	var ttransport thrift.TTransport
	var protocolFactory *thrift.TBinaryProtocolFactory
	var err error

	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory = thrift.NewTTransportFactory()
	ttransport, err = thrift.NewTSocket(address)
	ttransport = transportFactory.GetTransport(ttransport)
	if err = ttransport.Open(); err != nil {
		//logger.Println("Failed to Open Transport", transport, protocolFactory)
		return nil, nil
	}
	return ttransport, protocolFactory
}
