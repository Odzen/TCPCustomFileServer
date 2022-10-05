package types

import "net"

type Message struct {
	Text    string
	Address string
	Channel string
}

func NewMessage(msg string, conn net.Conn, channel string) Message {
	addr := conn.RemoteAddr().String()
	return Message{
		Text:    addr + msg,
		Address: addr,
		Channel: channel,
	}
}
