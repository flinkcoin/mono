package voting

import (
	"log"
	"sync"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// Placeholder structs and interfaces
type ByteString []byte
type FullBlock struct {
	Block *Block
}
type Block struct {
	BlockHash *Hash
	Body      *BlockBody
}
type BlockBody struct {
	AccountId ByteString
}
type Hash struct {
	Hash ByteString
}
type BlockVerifyStock interface {
	GetBlock(blockHash ByteString) []FullBlock
	Remove(blockHash ByteString)
}
type Storage interface {
	NewTransaction(func(t fdb.Transaction) error) error
}
type AccountCache interface {
	SetLastBlockHash(accountId, blockHash ByteString)
}

type BlockVerifyVotingService struct {
	blockVerifyStock BlockVerifyStock
	storage          Storage
	accountCache     AccountCache
	publishProcessor chan ByteString
	mu               sync.Mutex
}

func NewBlockVerifyVotingService(blockVerifyStock BlockVerifyStock, storage Storage, accountCache AccountCache) *BlockVerifyVotingService {
	service := &BlockVerifyVotingService{
		blockVerifyStock: blockVerifyStock,
		storage:          storage,
		accountCache:     accountCache,
		publishProcessor: make(chan ByteString, 1000),
	}
	go service.startProcessor()
	return service
}

func (s *BlockVerifyVotingService) startProcessor() {
	for blockHash := range s.publishProcessor {
		s.process(blockHash)
	}
}

func (s *BlockVerifyVotingService) NewBlock(blockHash ByteString) {
	s.publishProcessor <- blockHash
}

func (s *BlockVerifyVotingService) process(blockHash ByteString) {
	blocks := s.blockVerifyStock.GetBlock(blockHash)

	if len(blocks) == 0 {
		log.Println("Something not ok here, block missing!")
		return
	}

	err := s.storage.NewTransaction(func(t fdb.Transaction) error {
		for _, fb := range blocks {
			s.persistBlock(fb, t)
		}
		return nil
	})
	if err != nil {
		log.Printf("Could not write vote result to DB: %v", err)
		return
	}

	for _, fb := range blocks {
		body := fb.Block.Body
		block := fb.Block
		s.accountCache.SetLastBlockHash(body.AccountId, block.BlockHash.Hash)
	}

	s.blockVerifyStock.Remove(blockHash)
}

func (s *BlockVerifyVotingService) persistBlock(fb FullBlock, t fdb.Transaction) {
	// Implement the logic to persist the block using the transaction
}
