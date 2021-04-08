package rabbitmq

import (
	"context"
	"math/rand"
	"rpc"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

type opt struct {
}

type rabbit struct {
	opt  rpc.Options
	log  rpc.ILogger
	conn *amqp.Connection
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// Publish
// The function Publish is a function calling rpc server .
// The request param is a struct that implements  rpc.RequestMsgRPC
func (this rabbit) Publish(ctx context.Context, request rpc.RequestMsgRPC) (err error) {
	ch, err := this.conn.Channel()

	if err != nil {
		this.log.Errorw("err-msg-rabbit", "err", err.Error())
		return
	}
	defer ch.Close()

	msg_proto, err_encode := proto.Marshal(request.GetMsg())

	if err_encode != nil {
		this.log.Errorw("err-msg-rabbit", "err", err.Error())
		return
	}
	corrId := randomString(32)

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	publish_data := amqp.Publishing{
		ContentType:   "text/plain",
		CorrelationId: corrId,
		Body:          msg_proto,
		ReplyTo:       q.Name,
	}

	if err != nil {
		this.log.Errorw("err-msg-rabbit", "err", err.Error())
		return
	}

	err = ch.Publish(
		request.GetExchange(),   // exchange
		request.GetRoutingKey(), // routing key
		false,                   // mandatory
		false,                   // immediate
		publish_data,
	)

	if err != nil {
		this.log.Errorw("err-msg-rabbit", "err", err.Error())
		return
	}

	if request.HaveReply() {
		msgs, _ := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		if err != nil {
			this.log.Errorw("err-msg-rabbit", "err", err.Error())
			return
		}

		select {
		case d, ok := <-msgs:
			this.log.Infow("msg", "msg", d)
			if ok {
				proto.Unmarshal(d.Body, request.GetReplyMsg())
				return
			}
		case <-time.After(time.Minute * 10):
			this.log.Errorw("err-timeout-reply-rabbit")
			return
		}
	}

	return
}

// Event server
// To setup the server listen queue , You have to set 3 Properties those are exchange , routing key , function executed
// The RequestRPC param is struct that implements the RequestRPC Interface
func (this rabbit) EventServer(ctx context.Context, request rpc.RequestRPC) (err error) {
	ch, err := this.conn.Channel()
	this.log.Infow("starting EventServer")
	if err != nil {
		this.log.Errorw("err-msg-rabbit", "err", err.Error())
		return
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		this.log.Infow("err-msg-rabbit", "err", err.Error())
		return
	}

	err = ch.ExchangeDeclare(
		request.GetExchange(), // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)

	if err != nil {
		this.log.Infow("err-msg-rabbit", "err", err.Error())
		return
	}
	err = ch.QueueBind(
		q.Name,
		request.GetExchange(),
		request.GetRoutingKey(),
		false,
		nil,
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		this.log.Infow("err-msg-rabbit", "err", err.Error())
		return
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				this.log.Infow("Recovered in f", "recovery", r)
			}
		}()

		for d := range msgs {
			this.log.Infow("msg-rabbit", "msg", d)
			var msg proto.Message

			msg, err = request.GetRpcServer()(ctx)

			msg_proto, err_encode := proto.Marshal(msg)

			if err_encode == nil && d.ReplyTo != "" && d.CorrelationId != "" {
				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          msg_proto,
					})
			}

			d.Ack(false)
		}
	}()

	select {
	case <-ctx.Done():
		return
	}
}

func Instance(ctx context.Context, opt rpc.Options, logger rpc.ILogger) rabbit {
	conn, err := amqp.Dial(opt.URIConnection())
	if err != nil {
		logger.Infow("err-msg-rabbit", "err", err.Error())
	}
	return rabbit{opt, logger, conn}
}
