package report

import (
	"errors"
	"net"

	"git.apache.org/thrift.git/lib/go/thrift"
	"git.vodjk.com/golang/common/modules/net/flume"
)

// Flume 通信结构体
type Flume struct {
	service *flume.ThriftSourceProtocolClient
	headers map[string]string
}

// NewFlume 创建一个Flume连接
func NewFlume(host string, port string) (*Flume, error) {
	if len(host) == 0 {
		return nil, errors.New("The host option is required.")
	}
	if len(port) == 0 {
		return nil, errors.New("The port option is required.")
	}

	f := new(Flume)
	var trans thrift.TTransport
	trans, err := thrift.NewTSocket(net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}

	// newTransport
	trans = thrift.NewTFramedTransport(trans)
	f.service = flume.NewThriftSourceProtocolClientFactory(trans, thrift.NewTCompactProtocolFactory())
	f.headers = make(map[string]string)

	return f, nil
}

func (f *Flume) Write(p []byte) (int, error) {
	event := &flume.ThriftFlumeEvent{
		Headers: f.headers,
		Body: p,
	}

	if !f.service.Transport.IsOpen() {
		f.service.Transport.Open()
	}

	_, err := f.service.Append(event)
	if err != nil {
		f.service.Transport.Close()
	}
	return 0, nil
}
