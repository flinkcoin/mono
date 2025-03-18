package messaging

import (
	"github.com/flinkcoin/mono/libs/shared/pkg/base"
	nats "github.com/nats-io/nats.go"
)

type Queue struct {
	conn *nats.Conn
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Connect() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		base.Log.Error("Error connecting to NATS: %v:", err)
	}
	q.conn = nc
}

func (q *Queue) Close() {
	q.conn.Close()
}

func (q *Queue) Publish(subject string, data []byte) {
	err := q.conn.Publish(subject, data)
	if err != nil {
		base.Log.Error("Error publishing message: %v:", err)
	}
}

func (q *Queue) Subscribe() {
	q.conn.Subscribe("subject", func(msg *nats.Msg) {
		base.Log.Info("Received message: %s", string(msg.Data))
	})
}
