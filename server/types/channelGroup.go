package types

import "fmt"

type channel map[int][]Client

type ChannelGroup struct {
	Channels channel
}

func NewChannelGroup(channels channel) ChannelGroup {
	return ChannelGroup{
		Channels: channels,
	}
}

func (channelGroup *ChannelGroup) GetClientsByChannel(channel int) []Client {
	return channelGroup.Channels[channel]
}

func (channelGroup *ChannelGroup) SuscribeToChannelGroup(client Client, channel int) {
	channelGroup.Channels[channel] = append(channelGroup.Channels[channel], client)
}

func (channelGroup *ChannelGroup) DeleteClientFromChannel(client Client, channel int) {
	clients := channelGroup.Channels[channel]
	indexOfClient := channelGroup.getIndexClientFromChannel(client, channel)
	clientsAfterRemoval := append(clients[:indexOfClient], clients[indexOfClient+1:]...)
	channelGroup.Channels[channel] = clientsAfterRemoval
}

func (channelGroup *ChannelGroup) getIndexClientFromChannel(wantedClient Client, channel int) int {
	clientsInChannel := channelGroup.Channels[channel]
	for index, client := range clientsInChannel {
		if client.equals(wantedClient) {
			return index
		}
	}
	return -1
}

func (channelGroup ChannelGroup) Print() {
	for channel, clients := range channelGroup.Channels {
		fmt.Println("Channel: ", channel, " --> clients: ", clients)
	}
}
