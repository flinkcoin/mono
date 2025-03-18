package process

import (
	"fmt"
	"github.com/flinkcoin/mono/libs/schema/pkg/broker"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

// Process holds a map of handlers
type Process struct {
	handlers map[string]func(*anypb.Any, []byte) error
}

// NewProcess creates a new Process and initializes the handlers
func NewProcess() *Process {
	p := &Process{
		handlers: make(map[string]func(*anypb.Any, []byte) error),
	}
	p.initHandlers()
	return p
}

// initHandlers registers functions to handle various messages
func (p *Process) initHandlers() {
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.PaymentReq"] = p.handlePaymentReq
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockPub"] = p.handleBlockPub
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockConfirmPub"] = p.handleBlockConfirmPub
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockVerifyPub"] = p.handleBlockVerifyPub
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockVerifyConfirmPub"] = p.handleBlockVerifyConfirmPub
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockReq"] = p.handleBlockReq
	p.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockRes"] = p.handleBlockRes
}

// ProcessMsg unmarshals the message and delegates to the correct handler
func (p *Process) ProcessMsg(payload []byte) error {
	var msg broker.Message

	err := proto.Unmarshal(payload, &msg)
	if err != nil {
		log.Fatalf("Failed to deserialize message: %v", err)
	}

	if msg.Any == nil {
		return fmt.Errorf("any is nil")
	}

	handler, found := p.handlers[msg.Any.TypeUrl]
	if !found {
		return fmt.Errorf("unknown type URL: %s", msg.Any.TypeUrl)
	}

	return handler(msg.Any, payload)
}

func (p *Process) handlePaymentReq(any *anypb.Any, payload []byte) error {
	return nil
}

func (p *Process) handleBlockPub(any *anypb.Any, payload []byte) error {
	var blockPub broker.Message_BlockPub
	if err := anypb.UnmarshalTo(any, &blockPub, proto.UnmarshalOptions{}); err != nil {
		return err
	}

	body := blockPub.Body

	if body.MsgId==nil{
		return fmt.Errorf("msgId is nil")
	}
	else{
		"Add to block is queue, to do!"
	}





	if !p.cryptoManager.VerifyData(body.NodeId, []byte(fmt.Sprintf("%v", body)), blockPub.Signature) {
		return fmt.Errorf("failed to verify signature")
	}




}

func (p *Process) handleBlockConfirmPub(any *anypb.Any, payload []byte) error {
	var blockConfirm broker.Message_BlockConfirmPub
	if err := anypb.UnmarshalTo(any, &blockConfirm, proto.UnmarshalOptions{}); err != nil {
		return err
	}

	body := blockConfirm.Body

	if !p.cryptoManager.VerifyData(body.NodeId, []byte(fmt.Sprintf("%v", body)), blockConfirm.Signature) {
		return fmt.Errorf("failed to verify signature")
	}

}

func (p *Process) handleBlockVerifyPub(any *anypb.Any, payload []byte) error {
	var verifyPub broker.Message_BlockVerifyPub
	if err := anypb.UnmarshalTo(any, &verifyPub, proto.UnmarshalOptions{}); err != nil {
		return err
	}
}

func (p *Process) handleBlockVerifyConfirmPub(any *anypb.Any, payload []byte) error {
	var verifyConfirmPub broker.Message_BlockVerifyConfirmPub
	if err := anypb.UnmarshalTo(any, &verifyConfirmPub, proto.UnmarshalOptions{}); err != nil {
		return err
	}
}

func (p *Process) handleBlockReq(any *anypb.Any, payload []byte) error {
	var blockReq broker.Message_BlockReq
	if err := anypb.UnmarshalTo(any, &blockReq, proto.UnmarshalOptions{}); err != nil {
		return err
	}
}

func (p *Process) handleBlockRes(any *anypb.Any, payload []byte) error {
	var blockRes broker.Message_BlockRes
	if err := anypb.UnmarshalTo(any, &blockRes, proto.UnmarshalOptions{}); err != nil {
		return err
	}
}
