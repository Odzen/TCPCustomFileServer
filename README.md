# TCPCustomFileServer
Create a server that allows you to transfer files between 2 or more clients using a custom protocol (non-standard protocol) based on TCP.

# Commands for clients
- `=username <name>` - set an username. If the clients connects and doesn't set an username, the connection will stay anonymous.
- `=suscribe <number of the channel>` - join a channel. If the channel doesn't exist, it will be created. The client can be in one channel at the same time. The channel must be a number
- `=channels` -  show a list of available channels to join
- `=current` -  shows the channel to which the client is subscribed
- `=intructions` -  shows to the clients the available commands and their functionalities
- `=message <string message>` - broadcast message to everyone in the channel
- `=file <file>` - broadcast file to every client in the channel
- `=exit` - disconnects from the server.
