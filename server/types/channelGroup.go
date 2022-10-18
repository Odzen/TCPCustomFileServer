package types

import (
	"fmt"
	"log"
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
	if client.suscribedToChannel != 0 {
		fmt.Fprintln(client.Connection, "Client suscribed to another channel: ")
		channelGroup.DeleteClientFromChannel(*client, client.suscribedToChannel)
	}

	_, found := channelGroup.Channels[channel]

	// If the channel doesn't exists, create the channel and add the client to the channel
	if !found {
		channelGroup.Channels[channel] = []*Client{client}
		client.suscribedToChannel = channel
		return
	}

	// If the channel exist, just add the client to the channel
	channelGroup.Channels[channel] = append(channelGroup.Channels[channel], client)
	client.suscribedToChannel = channel

	// Notify clients
	channelGroup.Broadcast(NewMessage(fmt.Sprintf("with the name: %s, has joined the room.", client.Name), client.Connection, channel))
	fmt.Fprintln(client.Connection, "-> "+"Welcome to the channel # "+strconv.Itoa(channel))

}

func (channelGroup *ChannelGroup) DeleteClientFromChannel(client Client, channel int) {
	clients := channelGroup.Channels[channel]
	indexOfClient := channelGroup.getIndexClientFromChannel(&client, client.suscribedToChannel)
	clientsAfterRemoval := append(clients[:indexOfClient], clients[indexOfClient+1:]...)
	channelGroup.Channels[channel] = clientsAfterRemoval
	client.suscribedToChannel = 0
	channelGroup.Broadcast(NewMessage(" has left the channel.", client.Connection, channel))
}

func (channelGroup *ChannelGroup) getIndexClientFromChannel(wantedClient *Client, channel int) int {
	clientsInChannel := channelGroup.Channels[channel]
	for index, client := range clientsInChannel {
		if client.equals(*wantedClient) {
			log.Println("Equal", client)
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
