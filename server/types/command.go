package types

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Odzen/TCPCustomFileServer/utils"
)

type idCommand int

const (
	USERNAME idCommand = iota
	SUBSCRIBE
	CHANNELS
	MESSAGE
	CURRENT_CHANNEL
	INSTRUCTIONS
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
	case "=subscribe":
		client.Commands <- Command{
			Id:     SUBSCRIBE,
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
	case "=current":
		client.Commands <- Command{
			Id:     CURRENT_CHANNEL,
			Client: client,
			Args:   args,
		}
	case "=instructions":
		client.Commands <- Command{
			Id:     INSTRUCTIONS,
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
		fmt.Fprintf(client.Connection, "-> The command `%s` was not accepted. Use the command `=instructions` to see the available commands \n", command)
	}
}

func CreateUsername(client *Client, args []string) {
	client.ChangeName(args[1])
	fmt.Fprintln(client.Connection, "-> Your username has been changed to: "+client.Name)
}

func SuscribeToChannel(client *Client, args []string, channelGroup ChannelGroup) {
	selectedChannel, err := strconv.Atoi(args[1])

	if err != nil {
		fmt.Fprintln(client.Connection, "-> The Channel must be a number!")
		return
	}

	channelGroup.SuscribeToChannelGroup(client, selectedChannel)

	channelGroup.Print()
}

func ShowChannels(client *Client, args []string, channelGroup ChannelGroup) {
	fmt.Fprintf(client.Connection, "-> Available channels: %v \n", channelGroup.GetAvailableChannels())
}

func CurrentChannel(client *Client) {
	if client.SuscribedToChannel == 0 {
		fmt.Fprintf(client.Connection, "-> You're not subscribed to any channel, use the command ´=channels´ to see available channels or create a new one by using the command ´=subscribe <number>´\n")
		return
	}
	fmt.Fprintf(client.Connection, "-> You're subscribed to the channel # %s \n", strconv.Itoa(client.GetCurrentChannel()))
}

func Instructions(client *Client) {
	fmt.Fprintf(client.Connection, "-> `=username <name>` \n-> `=suscribe <number of the channel>` \n-> `=channels` \n-> `=current` \n-> `=intructions` \n-> `=message <string message>` \n-> `=file <file>`\n-> `=exit \n")
}

func SendMessage(client *Client, args []string, channelGroup ChannelGroup) {
	if client.SuscribedToChannel == 0 {
		fmt.Fprintln(client.Connection, "-> Subscribe first to a channel to send messages")
		return
	}

	channelGroup.BroadcastMessage(NewMessage(fmt.Sprintln("--"+client.Name+"-- : "+strings.Join(args[1:], " ")), client.Connection, client.SuscribedToChannel))
}

func SendFile(client *Client, args []string, channelGroup ChannelGroup) {

	if client.SuscribedToChannel == 0 {
		fmt.Fprintln(client.Connection, "-> Subscribe first to a channel to send files")
		return
	}

	fileToSend, err := ProccessingFile(client.Connection, args[1], client)

	if err {
		log.Println("Error processing file")
		return
	}

	channelGroup.BroadcastFile(fileToSend)

	// fileToSend.SendFileToClient()

}

func Exit(client *Client, channelGroup ChannelGroup) {
	log.Printf("Client left: %s", client.Address)

	channelGroup.Print()
	if client.SuscribedToChannel != 0 {
		channelGroup.DeleteClientFromChannel(*client, client.SuscribedToChannel)
	}

	utils.CloseConnectionClient(client.Connection)
}
