package types

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type idCommand int

const (
	USERNAME idCommand = iota
	SUSCRIBE
	CHANNELS
	MESSAGE
	FILE
	EXIT
)

type Command struct {
	Id     idCommand
	Client *Client
	Args   []string
}

func ProcessCommand(command string, args []string, client *Client) {
	switch command {
	case "=username":
		client.Commands <- Command{
			Id:     USERNAME,
			Client: client,
			Args:   args,
		}
	case "=suscribe":
		client.Commands <- Command{
			Id:     SUSCRIBE,
			Client: client,
			Args:   args,
		}
	case "=channels":
		client.Commands <- Command{
			Id:     CHANNELS,
			Client: client,
			Args:   args,
		}
	case "=message":
		client.Commands <- Command{
			Id:     MESSAGE,
			Client: client,
			Args:   args,
		}
	case "=file":
		client.Commands <- Command{
			Id:     FILE,
			Client: client,
			Args:   args,
		}
	case "=exit":
		client.Commands <- Command{
			Id:     EXIT,
			Client: client,
			Args:   args,
		}
	default:
		fmt.Fprintln(client.Connection, "Unknown Command: "+command)
	}
}

func CreateUsername(client *Client, args []string) {
	client.ChangeName(args[1])
	fmt.Fprintln(client.Connection, "Username has been changed to: "+client.Name)
}

func SuscribeToChannel(client *Client, args []string, channelGroup ChannelGroup) {
	selectedChannel, err := strconv.Atoi(args[1])

	if err != nil {
		fmt.Fprintln(client.Connection, "The Channel must be a number!")
		return
	}

	channelGroup.SuscribeToChannelGroup(client, selectedChannel)

	channelGroup.Print()
}

func ShowChannels(client *Client, args []string, channelGroup ChannelGroup) {
	fmt.Fprintf(client.Connection, "Available channels: %v \n", channelGroup.GetAvailableChannels())
}

func SendMessage(client *Client, args []string, channelGroup ChannelGroup) {
	if client.suscribedToChannel == 0 {
		fmt.Fprintln(client.Connection, "Suscribe to a channel to send messages")
		return
	}

	channelGroup.Broadcast(NewMessage(fmt.Sprintln("--"+client.Name+"-- texted : "+strings.Join(args[1:], " ")), client.Connection, client.suscribedToChannel))
}

func SendFile(client *Client, args []string) {

}

func Exit(client *Client, channelGroup ChannelGroup) {
	log.Printf("Client left: %s", client.Address)
	channelGroup.Print()
	if client.suscribedToChannel != 0 {
		channelGroup.DeleteClientFromChannel(*client, client.suscribedToChannel)
	}

}
