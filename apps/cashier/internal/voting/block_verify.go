package voting

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"log"
	"sync"
)

type BlockVerify struct {
	commonHandler    CommonProcessor
	nodeManager      NodeManager
	cryptoManager    CryptoManager
	blockCache       BlockCache
	blockVerifyStock BlockVerifyStock
	publishProcessor chan ByteString
	mu               sync.Mutex
}

func NewBlockVerifyService(commonHandler CommonProcessor, nodeManager NodeManager, cryptoManager CryptoManager, blockCache BlockCache, blockVerifyStock BlockVerifyStock) *BlockVerify {
	service := &BlockVerify{
		commonHandler:    commonHandler,
		nodeManager:      nodeManager,
		cryptoManager:    cryptoManager,
		blockCache:       blockCache,
		blockVerifyStock: blockVerifyStock,
		publishProcessor: make(chan ByteString, 1000),
	}
	go service.startProcessor()
	return service
}

func (s *BlockVerify) startProcessor() {
	for blockHash := range s.publishProcessor {
		s.process(blockHash)
	}
}

func (s *BlockVerify) NewBlock(blockHash ByteString) {
	s.publishProcessor <- blockHash
}

func (s *BlockVerify) process(blockHash ByteString) {
	if !s.validateBlock(blockHash) {
		return
	}

	block, found := s.blockCache.GetBlock(blockHash)
	if !found {
		return
	}

	body := &BlockConfirmPub_Body{
		BlockHash: block.Block.BlockHash.Hash,
		MsgId:     ByteString(uuid.New().NodeID()),
		NodeId:    s.nodeManager.GetNodeId(),
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

	blockExistConfirmPub := &BlockConfirmPub{
		Body:      body,
		Signature: signature,
	}

	anyMsg, err := ptypes.MarshalAny(blockExistConfirmPub)
	if err != nil {
		log.Printf("Error creating Any message: %v", err)
		return
	}

	err = s.commonHandler.Flood(anyMsg)
	if err != nil {
		log.Printf("Error flooding message: %v", err)
		return
	}
}

func (s *BlockVerify) validateBlock(blockHash ByteString) bool {
	// Implement block validation logic here
	return true
}
