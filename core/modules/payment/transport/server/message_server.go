package server

import (
	"fmt"

	paymentmessaging "wechat-clone/core/modules/payment/application/messaging"
	"wechat-clone/core/shared/pkg/stackErr"
)

type Server interface {
	Start() error
	Stop() error
}

type messageServer struct {
	messageHandler paymentmessaging.MessageHandler
}

func NewMessageServer(messageHandler paymentmessaging.MessageHandler) (Server, error) {
	if messageHandler == nil {
		return nil, stackErr.Error(fmt.Errorf("message handler can not be nil"))
	}
	return &messageServer{messageHandler: messageHandler}, nil
}

func (s *messageServer) Start() error {
	return s.messageHandler.Start()
}

func (s *messageServer) Stop() error {
	return s.messageHandler.Stop()
}
