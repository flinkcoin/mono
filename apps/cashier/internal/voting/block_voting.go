package voting

import (
	"log"
	"sync"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

type BlockVoting struct {
	blockStock         BlockStock
	storage            Storage
	accountCache       AccountCache
	accountServiceImpl AccountServiceImpl
	publishProcessor   chan ByteString
	mu                 sync.Mutex
}

func NewBlockVoting(blockStock BlockStock, storage Storage, accountCache AccountCache, accountServiceImpl AccountServiceImpl) *BlockVoting {
	service := &BlockVoting{
		blockStock:         blockStock,
		storage:            storage,
		accountCache:       accountCache,
		accountServiceImpl: accountServiceImpl,
		publishProcessor:   make(chan ByteString, 1000),
	}
	go service.startProcessor()
	return service
}

func (s *BlockVoting) startProcessor() {
	for blockHash := range s.publishProcessor {
		s.process(blockHash)
	}
}

func (s *BlockVoting) NewBlock(blockHash ByteString) {
	s.publishProcessor <- blockHash
}

func (s *BlockVoting) process(blockHash ByteString) {
	block, found := s.blockStock.GetBlock(blockHash)

	if !found {
		log.Println("Something not ok here, block missing!")
		return
	}

	err := s.storage.NewTransaction(func(t fdb.Transaction) error {
		previousBlock, found := s.storage.GetBlock(t, block.Body.PreviousBlock)
		if found {
			previousBlock.Block.Body.NextBlock = block.BlockHash.Hash
			s.storage.PutBlock(t, previousBlock.Block.BlockHash.Hash, previousBlock)
		}

		fullBlock := &FullBlock{Block: block}
		s.persistBlock(fullBlock, t)
		return nil
	})
	if err != nil {
		log.Printf("Could not write vote result to DB: %v", err)
		return
	}

	s.accountCache.SetLastBlockHash(block.Body.AccountId, block.BlockHash.Hash)
	s.blockStock.Remove(blockHash)

	infoRes := &ApiInfoRes{
		InfoType:     BlockConfirmType,
		BlockConfirm: &BlockConfirm{BlockHash: blockHash},
		AccountId:    block.Body.AccountId,
	}
	s.accountServiceImpl.SendInfo(infoRes)

	if block.Body.BlockType == SendType {
		infoRes = &ApiInfoRes{
			InfoType:        PaymentReceivedType,
			PaymentReceived: &PaymentReceived{BlockHash: blockHash},
			AccountId:       block.Body.SendAccountId,
		}
		s.accountServiceImpl.SendInfo(infoRes)
	}
}

func (s *BlockVoting) persistBlock(fb *FullBlock, t fdb.Transaction) {
	// Implement the logic to persist the block using the transaction
}
