package voting

import (
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
)

// Placeholder structs and interfaces
type ByteString []byte
type Node struct{}
type NodeAddress struct{}
type Pair struct {
	First  *Node
	Second *NodeAddress
}
type CommonProcessor interface {
	Flood(msg *any.Any) error
}
type NodeManager interface {
	GetNodeId() ByteString
}
type CryptoManager interface {
	SignData(data []byte) (ByteString, error)
}
type NodeVoting interface {
	NewNodeVote(pair Pair)
}

type NodeService struct {
	commonHandler    CommonProcessor
	nodeManager      NodeManager
	cryptoManager    CryptoManager
	nodeVoting       NodeVoting
	publishProcessor chan Pair
	mu               sync.Mutex
}

func NewNodeService(commonHandler CommonProcessor, nodeManager NodeManager, cryptoManager CryptoManager, nodeVoting NodeVoting) *NodeService {
	service := &NodeService{
		commonHandler:    commonHandler,
		nodeManager:      nodeManager,
		cryptoManager:    cryptoManager,
		nodeVoting:       nodeVoting,
		publishProcessor: make(chan Pair, 1000),
	}
	go service.startProcessor()
	return service
}

func (s *NodeService) startProcessor() {
	for pair := range s.publishProcessor {
		s.process(pair)
	}
}

func (s *NodeService) NewNode(pair Pair) {
	s.publishProcessor <- pair
}

func (s *NodeService) process(pair Pair) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node := pair.First
	nodeAddress := pair.Second

	body := &NodePub_Body{
		NodeId:      s.nodeManager.GetNodeId(),
		MsgId:       ByteString(uuid.New().NodeID()),
		Node:        node,
		NodeAddress: nodeAddress,
	}

	bodyBytes, err := proto.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling body: %v", err)
		return
	}

	signature, err := s.cryptoManager.SignData(bodyBytes)
	if err != nil {
		log.Printf("Error signing data: %v", err)
		return
	}

	nodeConfirmPub := &NodePub{
		Body:      body,
		Signature: signature,
	}

	anyMsg, err := ptypes.MarshalAny(nodeConfirmPub)
	if err != nil {
		log.Printf("Error creating Any message: %v", err)
		return
	}

	err = s.commonHandler.Flood(anyMsg)
	if err != nil {
		log.Printf("Error flooding message: %v", err)
		return
	}

	s.nodeVoting.NewNodeVote(Pair{First: node, Second: nodeAddress})
}
