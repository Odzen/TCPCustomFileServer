package types

type ChannelGroup map[string][]Client

func (channelGroup ChannelGroup) SuscribeToChannelGroup(client Client, channel string) {
	channelGroup[channel] = append(channelGroup[channel], client)
}
