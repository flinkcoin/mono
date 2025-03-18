package messaging

import (
	"github.com/flinkcoin/mono/apps/cashier/internal/process"
	"github.com/flinkcoin/mono/libs/schema/pkg/broker"
	"github.com/flinkcoin/mono/libs/schema/pkg/cashier"
	"github.com/flinkcoin/mono/libs/shared/pkg/base"
	nats "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Queue struct {
	conn *nats.Conn
	proc *process.Process
}

func NewQueue(proc *process.Process) *Queue {
	return &Queue{proc: proc}
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
	q.conn.Subscribe("subject", func(nmsg *nats.Msg) {
		var msg cashier.Msg

		// Deserialize the data into the Message object
		err := proto.Unmarshal(nmsg.Data, &msg)
		if err != nil {
			base.Log.Error("Failed to deserialize message: %v", err)
		}

		if msg.Message == nil {
			return
		}

		q.proc.ProcessMsg(nmsg.Data)
	})
}
