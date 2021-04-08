package main

import (
	"context"
	"rpc"
	"rpc/logger"
	"rpc/rabbitmq"

	"github.com/golang/protobuf/proto"
)

func main() {
	ctx := context.Background()
	opt := rpc.RabbitURI("amqp://admin:admin@localhost:5672/")

	log, _ := logger.NewLogger("test", "production")
	var irpc rpc.Irpc
	irpc = rabbitmq.Instance(ctx, opt, log)

	go func() {
		irpc.EventServer(ctx, rabbitmq.RabbitEventServer{
			WhereExchange:   "test",
			WhereRoutingKey: "test",
			WhereFunction: func(context.Context) (o proto.Message, e error) {
				o = &rpc.PathogenDTO{
					Id:     "12",
					Name:   "1",
					Avatar: "2",
				}
				return
			},
		})
	}()
	select {}
}
