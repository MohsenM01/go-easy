package messaging

import "errors"

type MessageRouter struct {
	handlers map[string]MessageHandler
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		handlers: make(map[string]MessageHandler),
	}
}

func (r *MessageRouter) RegisterHandler(queueName string, handler MessageHandler) {
	r.handlers[queueName] = handler
}

func (r *MessageRouter) RouteMessage(queueName string, msg []byte) error {
	handler, exists := r.handlers[queueName]
	if !exists {
		return errors.New("No handler registered for queue: " + queueName)
	}
	handler.ProcessMessage(string(msg))
	return nil
}
