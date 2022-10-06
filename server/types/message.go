package types

import "net"

type Message struct {
	Text            string
	Address         string
	ChannelPipeline int
}

func NewMessage(msg string, conn net.Conn, channel int) Message {
	addr := conn.RemoteAddr().String()
	return Message{
		Text:            addr + msg,
		Address:         addr,
		ChannelPipeline: channel,
	}
}
