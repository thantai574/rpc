package main

import (
	"context"
	"fmt"
	"rpc"
	"rpc/logger"
	"rpc/rabbitmq"
)

func main() {
	ctx := context.Background()
	opt := rpc.RabbitURI("amqp://admin:admin@localhost:5672/")

	log, _ := logger.NewLogger("test", "production")
	var irpc rpc.Irpc

	irpc = rabbitmq.Instance(ctx, opt, log)

	var res = &rpc.PathogenDTO{}

	irpc.Publish(ctx, rabbitmq.RabbitMsg{
		WhereExchange:   "test",
		WhereRoutingKey: "test",
		Msg: &rpc.PathogenDTO{
			Id:     "1",
			Name:   "1",
			Avatar: "2",
		},
		ReplyMsg: res,
		Reply:    true,
	})

	fmt.Print(res)

}
