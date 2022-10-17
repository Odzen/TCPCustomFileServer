package types

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
