package rabbitmq

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type RabbitEventServer struct {
	WhereExchange   string
	WhereRoutingKey string
	WhereFunction   func(ctx context.Context) (msg proto.Message, err error)
}

func (this RabbitEventServer) GetExchange() string {
	return this.WhereExchange
}

func (this RabbitEventServer) GetRoutingKey() string {
	return this.WhereRoutingKey
}

func (this RabbitEventServer) GetRpcServer() func(ctx context.Context) (msg proto.Message, err error) {
	return this.WhereFunction
}

type RabbitMsg struct {
	WhereExchange   string
	WhereRoutingKey string
	Msg             proto.Message
	ReplyMsg        proto.Message // Pointer
	Reply           bool
}

func (this RabbitMsg) GetExchange() string {
	return this.WhereExchange
}

func (this RabbitMsg) GetRoutingKey() string {
	return this.WhereRoutingKey
}

func (this RabbitMsg) HaveReply() bool {
	return this.Reply
}

func (this RabbitMsg) GetMsg() proto.Message {
	return this.Msg
}
func (this RabbitMsg) GetReplyMsg() proto.Message {
	return this.ReplyMsg
}
