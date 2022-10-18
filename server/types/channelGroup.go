package types

import (
	"fmt"
	"strconv"
)

type channel map[int][]*Client

type ChannelGroup struct {
	Channels channel
}

func NewChannelGroup(channels channel) ChannelGroup {
	return ChannelGroup{
		Channels: channels,
	}
}

func (channelGroup *ChannelGroup) GetClientsByChannel(channel int) []*Client {
	return channelGroup.Channels[channel]
}

func (channelGroup *ChannelGroup) SuscribeToChannelGroup(client *Client, channel int) {
	// If the client was suscribed to another channel before, remove it from that channel
	if client.SuscribedToChannel != 0 {
		fmt.Fprintln(client.Connection, "-> You were removed from the channel # ", client.SuscribedToChannel)
		channelGroup.DeleteClientFromChannel(*client, client.SuscribedToChannel)
	}

	channelGroup.Channels[channel] = append(channelGroup.Channels[channel], client)
	client.SuscribedToChannel = channel

	// Notify clients
	channelGroup.Broadcast(NewMessage(fmt.Sprintf(" with the name: %s, has joined the room.", client.Name), client.Connection, channel))
	fmt.Fprintln(client.Connection, "-> "+"Welcome to the channel # "+strconv.Itoa(channel))

}

func (channelGroup *ChannelGroup) DeleteClientFromChannel(client Client, channel int) {
	clients := channelGroup.Channels[channel]
	indexOfClient := channelGroup.getIndexClientFromChannel(&client, client.SuscribedToChannel)
	clientsAfterRemoval := append(clients[:indexOfClient], clients[indexOfClient+1:]...)
	channelGroup.Channels[channel] = clientsAfterRemoval
	client.SuscribedToChannel = 0
	channelGroup.Broadcast(NewMessage(" has left the channel.", client.Connection, channel))
}

func (channelGroup *ChannelGroup) getIndexClientFromChannel(wantedClient *Client, channel int) int {
	clientsInChannel := channelGroup.Channels[channel]
	for index, client := range clientsInChannel {
		if client.equals(*wantedClient) {
			return index
		}
	}
	return -1
}

func (channelGroup *ChannelGroup) GetAvailableChannels() []int {
	var channels []int
	for key := range channelGroup.Channels {
		channels = append(channels, key)
	}
	return channels
}

func (channelGroup *ChannelGroup) Print() {
	for channel, clients := range channelGroup.Channels {
		fmt.Printf("Channel: %d --> Clients : %v \n", channel, clients)
	}
}

func (channelGroup *ChannelGroup) Broadcast(msg Message) {
	for _, client := range channelGroup.Channels[msg.ChannelPipeline] {
		if msg.Address != client.Address { // Send the message to all the clients, exluding the one who sent it
			fmt.Fprintln(client.Connection, "-> "+msg.Text)
		}
	}
}
