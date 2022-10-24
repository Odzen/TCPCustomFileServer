package types

import "net"

type Message struct {
	Text            string
	AddressClient   string
	ChannelPipeline int
}

func NewMessage(msg string, conn net.Conn, channel int) Message {
	addressClient := conn.RemoteAddr().String()
	return Message{
		Text:            msg,
		AddressClient:   addressClient,
		ChannelPipeline: channel,
	}
}
