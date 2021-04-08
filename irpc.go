package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type RequestRPC interface {
	GetExchange() string
	GetRoutingKey() string
	GetRpcServer() func(ctx context.Context) (msg proto.Message, err error)
}

type RequestMsgRPC interface {
	GetExchange() string
	GetRoutingKey() string
	GetMsg() proto.Message
	GetReplyMsg() proto.Message
	HaveReply() bool
}

type Irpc interface {
	EventServer(ctx context.Context, request RequestRPC) (err error)
	Publish(ctx context.Context, request RequestMsgRPC) (err error)
}
