package messaging

type MessageHandler interface {
	ProcessMessage(msg string)
}
