package client

import (
	"github.com/giongto35/cloud-game/v2/pkg/ipc"
	"github.com/giongto35/cloud-game/v2/pkg/logger"
	"github.com/giongto35/cloud-game/v2/pkg/network"
)

type (
	NetClient interface {
		Close()
		Id() network.Uid
	}
	RegionalClient interface {
		In(region string) bool
	}
)

type SocketClient struct {
	NetClient

	id   network.Uid
	tag  string
	wire *ipc.Client
	log  *logger.Logger
}

func New(conn *ipc.Client, tag string, log *logger.Logger) SocketClient {
	id := network.NewUid()
	l := log.Wrap(log.With().Str("c-uid", id.Short()).Str("c-tag", tag))
	return SocketClient{id: id, wire: conn, tag: tag, log: l}
}

func (c SocketClient) Id() network.Uid { return c.id }

func (c SocketClient) Send(t uint8, data interface{}) ([]byte, error) { return c.wire.Call(t, data) }

func (c SocketClient) SendPacket(packet ipc.OutPacket) error { return c.wire.SendPacket(packet) }

func (c SocketClient) SendAndForget(t uint8, data interface{}) error { return c.wire.Send(t, data) }

func (c SocketClient) OnPacket(fn func(p ipc.InPacket)) { c.wire.OnPacket = fn }

func (c SocketClient) GetLogger() *logger.Logger { return c.log }

func (c SocketClient) Listen() { <-c.wire.Conn.Done }

func (c SocketClient) Close() { c.wire.Close() }

func (c SocketClient) String() string { return c.tag + ":" + c.Id().Short() }