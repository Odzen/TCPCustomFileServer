package server

type ChannelGroup map[int][]Client

func (channelGroup ChannelGroup) SuscribeToChannelGroup(client Client, channel int) {
	channelGroup[channel] = append(channelGroup[channel], client)
}
