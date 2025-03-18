package api

import (
	"log"
	"sync"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
)

// Placeholder structs and interfaces
type ByteString []byte
type CommonProcessor interface {
	Flood(msg *any.Any) error
}
type NodeManager interface {
	GetNodeId() ByteString
}
type CryptoManager interface {
	SignData(data []byte) (ByteString, error)
}
type BlockCache interface {
	GetBlock(blockHash ByteString) (*FullBlock, bool)
}
type AccountCache interface {
	GetLastBlockHash(accountId ByteString) (ByteString, bool)
}
type AccountUnclaimedCache interface {
	GetLastBlockHash(accountId ByteString) (ByteString, bool)
}
type UnclaimedBlockCache interface {
	GetLastUnclaimedBlock(blockHash ByteString) (ByteString, bool)
}
type Storage interface {
	Count(columnFamily string) (int64, error)
}
type StreamObserver interface {
	OnNext(response interface{})
	OnCompleted()
	OnError(err error)
}
type FullBlock struct {
	Block *Block
}
type Block struct {
	Body *BlockBody
}
type BlockBody struct {
	PreviousBlockHash ByteString
}
type PaymentRequest struct{}
type InfoRes struct{}
type ApiInfoRes struct{}
type ApiPaymentTransactionReq struct{}
type ApiPaymentTransactionRes struct{}
type ApiTransactionReq struct{}
type ApiTransactionRes struct{}
type ApiInfoReq struct{}
type ApiInfoRes struct{}
type AccountCountReq struct{}
type AccountCountRes struct{}
type GetBlockReq struct{}
type GetBlockRes struct{}
type LastBlockReq struct{}
type LastBlockRes struct{}
type ListBlockReq struct{}
type ListBlockRes struct{}
type ListUnclaimedBlockReq struct{}
type ListUnclaimedBlockRes struct{}

type AccountServiceImpl struct {
	nodeManager           NodeManager
	cryptoManager         CryptoManager
	blockCache            BlockCache
	accountCache          AccountCache
	accountUnclaimedCache AccountUnclaimedCache
	unclaimedBlockCache   UnclaimedBlockCache
	storage               Storage
	infoObservers         map[int]StreamObserver
	mu                    sync.Mutex
}

func NewAccountServiceImpl(nodeManager NodeManager, cryptoManager CryptoManager, blockCache BlockCache, accountCache AccountCache, accountUnclaimedCache AccountUnclaimedCache, unclaimedBlockCache UnclaimedBlockCache, storage Storage) *AccountServiceImpl {
	return &AccountServiceImpl{
		nodeManager:           nodeManager,
		cryptoManager:         cryptoManager,
		blockCache:            blockCache,
		accountCache:          accountCache,
		accountUnclaimedCache: accountUnclaimedCache,
		unclaimedBlockCache:   unclaimedBlockCache,
		storage:               storage,
		infoObservers:         make(map[int]StreamObserver),
	}
}

func (s *AccountServiceImpl) NumAccounts(request *AccountCountReq, responseObserver StreamObserver) {
	count, err := s.storage.Count("ACCOUNT")
	if err != nil {
		log.Printf("Not ok: %v", err)
	}

	responseObserver.OnNext(&AccountCountRes{Count: count})
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) GetBlock(request *GetBlockReq, responseObserver StreamObserver) {
	block, found := s.blockCache.GetBlock(request.BlockHash)
	if found {
		responseObserver.OnNext(&GetBlockRes{Block: block.Block})
	} else {
		responseObserver.OnNext(&GetBlockRes{})
	}
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) LastBlock(request *LastBlockReq, responseObserver StreamObserver) {
	lastBlockHash, found := s.accountCache.GetLastBlockHash(request.AccountId)
	if found {
		block, found := s.blockCache.GetBlock(lastBlockHash)
		if found {
			responseObserver.OnNext(&LastBlockRes{Block: block.Block})
		} else {
			responseObserver.OnNext(&LastBlockRes{})
		}
	} else {
		responseObserver.OnNext(&LastBlockRes{})
	}
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) ListBlocks(request *ListBlockReq, responseObserver StreamObserver) {
	lastBlockHash, found := s.accountCache.GetLastBlockHash(request.AccountId)
	if found {
		var blocks []*Block
		block, found := s.blockCache.GetBlock(lastBlockHash)
		for i := 0; i < request.Num && found; i++ {
			blocks = append(blocks, block.Block)
			block, found = s.blockCache.GetBlock(block.Block.Body.PreviousBlockHash)
		}
		responseObserver.OnNext(&ListBlockRes{Blocks: blocks})
	} else {
		responseObserver.OnNext(&ListBlockRes{})
	}
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) ListUnclaimedBlocks(request *ListUnclaimedBlockReq, responseObserver StreamObserver) {
	lastBlockHash, found := s.accountUnclaimedCache.GetLastBlockHash(request.AccountId)
	if found {
		var blocks []*Block
		blockHashes := []ByteString{lastBlockHash}
		lastUnclaimedBlock, found := s.unclaimedBlockCache.GetLastUnclaimedBlock(lastBlockHash)
		count := 0
		for found && count < request.Num {
			blockHashes = append(blockHashes, lastUnclaimedBlock)
			lastUnclaimedBlock, found = s.unclaimedBlockCache.GetLastUnclaimedBlock(lastUnclaimedBlock)
			count++
		}
		for _, blockHash := range blockHashes {
			block, found := s.blockCache.GetBlock(blockHash)
			if found {
				blocks = append(blocks, block.Block)
			} else {
				break
			}
		}
		responseObserver.OnNext(&ListUnclaimedBlockRes{Blocks: blocks})
	} else {
		responseObserver.OnNext(&ListUnclaimedBlockRes{})
	}
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) PaymentRequest(request *ApiPaymentTransactionReq, responseObserver StreamObserver) {
	paymentRequest := request.PaymentRequest
	err := s.publishPaymentRequest(paymentRequest)
	if err != nil {
		log.Printf("Something is wrong: %v", err)
	}
	responseObserver.OnNext(&ApiPaymentTransactionRes{Success: err == nil})
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) Transaction(request *ApiTransactionReq, responseObserver StreamObserver) {
	block := request.Block
	err := s.publish(block)
	if err != nil {
		log.Printf("Something is wrong: %v", err)
	}
	responseObserver.OnNext(&ApiTransactionRes{Success: err == nil})
	responseObserver.OnCompleted()
}

func (s *AccountServiceImpl) ReceiveInfos(request *ApiInfoReq, responseObserver StreamObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	observer, found := s.infoObservers[request.Id]
	if found {
		observer.OnCompleted()
		delete(s.infoObservers, request.Id)
	}
	s.infoObservers[request.Id] = responseObserver
}

func (s *AccountServiceImpl) SentInfo(infoRes *ApiInfoRes) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, observer := range s.infoObservers {
		observer.OnNext(infoRes)
		observer.OnCompleted()
		delete(s.infoObservers, id)
	}
}

func (s *AccountServiceImpl) publishPaymentRequest(paymentRequest *PaymentRequest) error {
	msgId := ByteString(uuid.New().NodeID())
	body := &PaymentReqBody{
		PaymentRequest: paymentRequest,
		MsgId:          msgId,
		NodeId:         s.nodeManager.GetNodeId(),
	}
	bodyBytes, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	signature, err := s.cryptoManager.SignData(bodyBytes)
	if err != nil {
		return err
	}
	paymentReq := &PaymentReq{
		Body:      body,
		Signature: signature,
	}
	s.sentInfo(&InfoRes{
		AccountId:      paymentRequest.ToAccountId,
		InfoType:       InfoRes_PAYMENT_REQUEST,
		PaymentRequest: paymentRequest,
	})
	return s.commonHandler.Flood(ptypes.MarshalAny(paymentReq))
}

func (s *AccountServiceImpl) publish(block *Block) error {
	msgId := ByteString(uuid.New().NodeID())
	body := &BlockPubBody{
		Block:  block,
		MsgId:  msgId,
		NodeId: s.nodeManager.GetNodeId(),
	}
	bodyBytes, err := proto.Marshal(body)
	if err != nil {
		return err
	}
	signature, err := s.cryptoManager.SignData(bodyBytes)
	if err != nil {
		return err
	}
	blockPub := &BlockPub{
		Body:      body,
		Signature: signature,
	}
	s.blockService.NewBlock(Pair{First: s.nodeManager.GetNodeId(), Second: block})
	return s.commonHandler.Flood(ptypes.MarshalAny(blockPub))
}
