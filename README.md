# Go RPC server
#### Introduction 
RPC server implements server rabbitmq 

##### Installation 
```shell
go get -u github.com/thantai574/rpc
```
##### Import 
```go
import rpc "github.com/thantai574/rpc"
```
##### Example Server RPC rabbitmq
```go
package main

import (
	"context"
	"github.com/thantai574/rpc"
	"github.com/thantai574/rpc/example/rabbit/msg"
	"github.com/thantai574/rpc/logger"
	"github.com/thantai574/rpc/rabbitmq"

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
			Function: func(context.Context) (o proto.Message, e error) {
				o = &msg.PathogenDTO{
					Id:     "1231212",
					Name:   "1",
					Avatar: "2",
				}
				return
			},
		})
	}()
	select {}
}

```
##