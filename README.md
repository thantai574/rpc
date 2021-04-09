# The RPC server uses the request-reply design pattern
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
##### Example setups a server RPC that implements rabbitmq
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

##### Example Client RPC rabbitmq
```go
package main

import (
	"context"
	"fmt"

	"github.com/thantai574/rpc"
	"github.com/thantai574/rpc/example/rabbit/msg"
	"github.com/thantai574/rpc/logger"
	"github.com/thantai574/rpc/rabbitmq"
)

func main() {
	ctx := context.Background()
	opt := rpc.RabbitURI("amqp://admin:admin@localhost:5672/")

	log, _ := logger.NewLogger("test", "production")
	var irpc rpc.Irpc

	irpc = rabbitmq.Instance(ctx, opt, log)

	var res = &msg.PathogenDTO{}

	irpc.Publish(ctx, rabbitmq.RabbitMsg{
		WhereExchange:   "test",
		WhereRoutingKey: "test",
		Msg: &msg.PathogenDTO{
			Id:     "1",
			Name:   "1",
			Avatar: "2",
		},
		ReplyMsg: res, // ReplyMsg 
		Reply:    true, // to wait till has a msg 
		Timeout:  time.Minute,
	})

	fmt.Print(res)

}
```
##