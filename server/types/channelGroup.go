package types

import (
	"encoding/json"
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

func (channelGroup *ChannelGroup) suscribeToChannel(client *Client, channel int) {
	// If the client was subscribed to another channel before, remove it from that channel
	if client.SubscribedToChannel != 0 {
		channelGroup.deleteClientFromChannel(*client, client.SubscribedToChannel)
		fmt.Fprintln(client.Connection, "-> You were removed from the channel # ", client.SubscribedToChannel)
	}

	channelGroup.Channels[channel] = append(channelGroup.Channels[channel], client)
	client.SubscribedToChannel = channel

	// Notify clients
	channelGroup.broadcastMessage(NewMessage(fmt.Sprintf(" %s, has joined the room.", client.Name), client.Connection, channel))
	fmt.Fprintln(client.Connection, "-> "+"Welcome to the channel # "+strconv.Itoa(channel))

}

func (channelGroup *ChannelGroup) deleteClientFromChannel(client Client, channel int) {
	clients := channelGroup.Channels[channel]
	indexOfClient := channelGroup.getIndexClientFromChannel(&client, client.SubscribedToChannel)
	clientsAfterRemoval := append(clients[:indexOfClient], clients[indexOfClient+1:]...)
	channelGroup.Channels[channel] = clientsAfterRemoval
	client.SubscribedToChannel = 0
	channelGroup.broadcastMessage(NewMessage(fmt.Sprintln(client.Name+" has left the channel."), client.Connection, channel))
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

func (channelGroup *ChannelGroup) getAvailableChannels() []int {
	var channels []int
	for key := range channelGroup.Channels {
		channels = append(channels, key)
	}
	return channels
}

func (channelGroup *ChannelGroup) print() {
	for channel, clients := range channelGroup.Channels {
		fmt.Printf("Channel %d : \n", channel)
		for _, client := range clients {
			fmt.Printf("%s // ", client.Connection.RemoteAddr().String()+"--"+client.Name)
		}
		fmt.Printf("\n")
	}
}

func (channelGroup *ChannelGroup) ToJson() ([]byte, error) {
	var clientsJSON []*Client

	// To show an empty array in the JSON format when the channel is empty
	clientsJSON = make([]*Client, 0)

	for _, clients := range channelGroup.Channels {
		clientsJSON = append(clientsJSON, clients...)
	}
	return json.Marshal(clientsJSON)
}

func (channelGroup *ChannelGroup) broadcastMessage(msg Message) {
	for _, client := range channelGroup.Channels[msg.ChannelPipeline] {
		if msg.AddressClient != client.Address { // Send the message to all the clients, exluding the one who sent it
			fmt.Fprintln(client.Connection, "-> "+msg.Text)
		}
	}
}

func (channelGroup *ChannelGroup) broadcastFile(file File) {

	file.appendToSentFiles()

	fmt.Println("Broadcasting file...")
	for _, client := range channelGroup.Channels[file.ChannelPipeline] {
		if file.AddressClient != client.Address { // Send the file to all the clients, exluding the one who sent it
			fmt.Fprintln(client.Connection, "-> "+"Sending File..")

			err := client.saveFile(file)

			if err != nil {
				fmt.Println("Error saving the file for the client:", client.Name+"--"+client.Address)
				fmt.Fprintln(client.Connection, "-> "+"Error saving the file")
				return
			}

			fmt.Fprintln(client.Connection, "-> "+"The file was saved successfully")

		}
	}
}
