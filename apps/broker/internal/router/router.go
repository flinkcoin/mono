package router

import (
	"bytes"
	"fmt"
	"github.com/flinkcoin/mono/apps/broker/internal/messaging"
	"github.com/flinkcoin/mono/libs/schema/pkg/broker"
	"github.com/flinkcoin/mono/libs/schema/pkg/core"
	"github.com/flinkcoin/mono/libs/shared/pkg/topics"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

type Router struct {
	queue    *messaging.Queue
	handlers map[string]func(*broker.Message, peer.ID, peer.PubKey) error
}

func NewRouter(queue *messaging.Queue) *Router {
	r := &Router{
		queue:    queue,
		handlers: make(map[string]func(*broker.Message, peer.ID, peer.PubKey) error),
	}
	r.initHandlers()
	return r
}

func (r *Router) initHandlers() {
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.IAmAlive"] = r.handleIAmAlive
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.AuthenticationReq"] = r.handleAuthenticationReq
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.AuthenticationRes"] = r.handleAuthenticationRes
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.NodePub"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.PaymentReq"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockPub"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockConfirmPub"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockVerifyPub"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockVerifyConfirmPub"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockReq"] = r.handleForward
	r.handlers["type.googleapis.com/flinkcoin.broker.Message.BlockRes"] = r.handleForward
}

func (r *Router) ProcessMsg(payload []byte, id peer.ID, pubKey peer.PubKey) error {
	var msg broker.Message

	// Deserialize the data into the Message object
	err := proto.Unmarshal(payload, &msg)
	if err != nil {
		log.Fatalf("Failed to deserialize message: %v", err)
	}

	if msg.Any == nil {
		return fmt.Errorf("any is nil")
	}

	handler, found := r.handlers[msg.Any.TypeUrl]
	if !found {
		return fmt.Errorf("unknown type URL: %s", msg.Any.TypeUrl)
	}

	return handler(msg, id, pubKey)
}

func (r *Router) handleIAmAlive(msg *broker.Message, id peer.ID, pubKey peer.PubKey) error {
	var specificMsg broker.Message_IAmAlive
	err := anypb.UnmarshalTo(any, &specificMsg, proto.UnmarshalOptions{})
	if err != nil {
		return fmt.Errorf("failed to unpack Any field: %v", err)
	}
	return nil
}

func (r *Router) handleAuthenticationReq(msg *broker.Message, id peer.ID, pubKey peer.PubKey) error {
	var authReq broker.Message_AuthenticationReq
	if err := anypb.UnmarshalTo(any, &authReq, proto.UnmarshalOptions{}); err != nil {
		return err
	}

	// Build the response body
	node := &core.Node{
		NodeId:    p.nodeManager.GetNodeId(),
		PublicKey: p.nodeManager.GetPublicKey(),
	}
	address := &core.NodeAddress{
		Ip:   p.config.Ip(),
		Port: p.config.Port(),
	}
	body := broker.Message_AuthenticationRes_Body{
		Node:        node,
		NodeAddress: address,
		Token:       authReq.Token,
	}

	// Sign the data
	signature, err := p.cryptoManager.SignData([]byte(fmt.Sprintf("%v", body)))
	if err != nil {
		return err
	}

	// Build the response
	authRes := broker.Message_AuthenticationRes{
		Body:      body,
		Signature: signature,
	}

	// Simulate writing the response back
	fmt.Println("Writing response:", authRes)
	return nil
}

func (r *Router) handleAuthenticationRes(msg *broker.Message, id peer.ID, pubKey peer.PubKey) error {
	var authRes broker.Message_AuthenticationRes
	if err := anypb.UnmarshalTo(any, &authRes, proto.UnmarshalOptions{}); err != nil {
		return err
	}

	body := authRes.Body
	node := body.Node
	nodeAddress := body.NodeAddress

	log.Printf("Node %s sending auth res!", node.NodeId)

	if !bytes.Equal(channelData.getToken(), body.Token) {
		log.Printf("Node %s not sending correct token!", node.NodeId)
		p.errorCloseChannel(ctx)
		return nil
	}

	if !p.nodeManager.CheckNodeExists(node.NodeId) {
		log.Printf("Node %s not in DB!", node.NodeId)
		p.discoveryServiceProvider().NewNode(node, nodeAddress)
		p.errorCloseChannel(ctx)
		return nil
	}

	if !p.cryptoManager.VerifyData(node.NodeId, []byte(fmt.Sprintf("%v", body)), authenticationRes.Signature) {
		log.Printf("Node %s sending wrong signature!", node.NodeId)
		p.errorCloseChannel(ctx)
		return nil
	}

	log.Printf("Node %s authentication success!", node.NodeId)
	channelData.setNodeId(node.NodeId)
	channelData.authenticate()

	commonHandler := p.commonHandlerProvider()
	if !commonHandler.AddNode(node.NodeId, "channelId") {
		log.Printf("Duplicate channel detected for nodeId: %s", node.NodeId)
		// ctx.close() // Simulate closing the context
	}

	return nil
}

func (r *Router) handleForward(msg *broker.Message, id peer.ID, pubKey peer.PubKey) error {

	r.queue.Publish(topics.CashierInbound.String(), payload)
	return nil
}
