package rpc_test

import (
	"context"
	"rpc"
	"rpc/logger"
	"rpc/rabbitmq"
	"testing"

	"google.golang.org/protobuf/runtime/protoiface"
)

func TestServer(t *testing.T) {
	// run server
	ctx := context.Background()
	opt := rpc.RabbitURI("amqp://admin:admin@localhost:5672/")

	log, _ := logger.NewLogger("test", "production")
	var irpc rpc.Irpc
	irpc = rabbitmq.Instance(ctx, opt, log)

	go func() {
		err := irpc.EventServer(ctx, rabbitmq.RabbitEventServer{
			WhereExchange:   "test",
			WhereRoutingKey: "test",
			WhereFunction:   func(context.Context) (protoiface.MessageV1, error) { panic("not implemented") },
		})

		if err != nil {
			t.Errorf("err")
		}
	}()
	select {}
}
