# TCP Custom Protocol
The server allows transferring files between several clients using a custom protocol (non-standard protocol) based on TCP.

## To run the server
`go run main.go server`. First, run the server, then the client. The server shows logs to know what is happening between the clients.

## To run the client
`go run main.go client`. The client can send messages or files through a channel. In addition, once a client receives a file, he/she will save the file in a specific folder `(outFiles/channel/client's name)`. For instance, if the client `juan` recieves a file through channel 8, he will put it in `(outFiles/8/juan)` folder.

## Folders and Files
To test the server, you can use the files located in the folder `testFiles` or any other file. Just make sure as a client, that the file doesn't exceed the limit of bytes that the server can handle (`50.5 KB`), and remember always to type the right path. The folder `outFiles` will be created automatically by the server, once the first file in the session is sent, broadcasted and received.

## Commands for clients
- `=username <name>` - set a username. If the client connects and doesn't set a username, the connection will stay anonymous.
- `=suscribe <number of the channel>` - join a channel. If the channel doesn't exist, it will be created. The client can be in one channel at the same time. The channel must be a number
- `=channels` -  show a list of available channels to join
- `=current` -  shows the channel to which the client is subscribed
- `=instructions` -  shows to the clients the available commands and their functionalities
- `=message <string message>` - broadcast message to everyone in the channel
- `=file <file>` - broadcast file to every client on the channel
- `=exit` - disconnects from the server.

## Endpoints
The server also has two open endpoints to receive and handle HTTP requests. The first one `/clients` shows which clients are connected and subscribed to specific channels. The second one `/files` shows file delivery statistics. This means that this server can also work as an API for any view. The port for listening to HTTP requests has to be different from the port for listening to TCP connections.

## Web Interface - View
[Link to GitHub Repository](https://github.com/Odzen/tcp-view)

## Demo
[Video](https://www.loom.com/share/29555cff5649472fba556bb14b2593d2)

[Scenarios](https://drive.google.com/file/d/1RufyrF8xTB7rilSWPYvE0OB0lcSUFNdA/view?usp=sharing)
