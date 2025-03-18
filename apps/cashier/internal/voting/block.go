package voting

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"log"
)

type BlockService struct {
	commonHandler CommonProcessor
	nodeManager   NodeManager
	cryptoManager CryptoManager
	blockVoting   BlockVoting
	blockStock    BlockStock
	blockHandler  ValidationHandler
}

func (s *BlockService) Process(pair Pair) {
	nodeId := pair.First
	block := pair.Second

	if !s.blockHandler.ValidateBlock(block) {
		return
	}

	err := s.blockStock.PutBlock(block.BlockHash.Hash, block)
	if err != nil {
		log.Printf("Error storing block: %v", err)
		return
	}

	body := &BlockConfirmPub_Body{
		BlockHash: block.BlockHash.Hash,
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

	s.blockVoting.NewBlockVote(Pair{First: nodeId, Second: block.BlockHash.Hash})
}
