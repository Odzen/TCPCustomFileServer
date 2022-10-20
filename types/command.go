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

func ProcessCommand(command string, args []string) (Command, bool) {
	switch command {
	case "=username":
		return Command{
			Id:   USERNAME,
			Args: args,
		}, true
	case "=subscribe":
		return Command{
			Id:   SUBSCRIBE,
			Args: args,
		}, true
	case "=channels":
		return Command{
			Id:   CHANNELS,
			Args: args,
		}, true
	case "=message":
		return Command{
			Id:   MESSAGE,
			Args: args,
		}, true
	case "=current":
		return Command{
			Id:   CURRENT_CHANNEL,
			Args: args,
		}, true
	case "=instructions":
		return Command{
			Id:   INSTRUCTIONS,
			Args: args,
		}, true
	case "=file":
		return Command{
			Id:   FILE,
			Args: args,
		}, true
	case "=exit":
		return Command{
			Id:   EXIT,
			Args: args,
		}, true
	default:
		//fmt.Fprintf(client.Connection, "-> The command `%s` was not accepted. Use the command `=instructions` to see the available commands \n", command)
		return Command{}, false
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

	channelGroup.Broadcast(NewMessage(fmt.Sprintln("--"+client.Name+"-- : "+strings.Join(args[1:], " ")), client.Connection, client.SuscribedToChannel))
}

func SendFile(client *Client, args []string) {
	fileToSend := ProccessingFile(client.Connection, args[1])
	fileToSend.SendFileToClient()
}

func Exit(client *Client, channelGroup ChannelGroup) {
	log.Printf("Client left: %s", client.Address)

	channelGroup.Print()
	if client.SuscribedToChannel != 0 {
		channelGroup.DeleteClientFromChannel(*client, client.SuscribedToChannel)
	}

	utils.CloseConnectionClient(client.Connection)
}
