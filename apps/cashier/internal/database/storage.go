package database

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/golang/protobuf/proto"
)

type Storage struct {
	*StorageBase
}

func NewStorage() (*Storage, error) {
	base, err := NewStorageBase()
	if err != nil {
		return nil, err
	}
	return &Storage{base}, nil
}

func (s *Storage) PutAccount(t fdb.Transaction, accountId, blockHash []byte) error {

	t.Set(s.dirs[ACCOUNT.Name].Pack(tuple.Tuple{accountId}), blockHash)
	return nil
}

func (s *Storage) GetAccount(t fdb.Transaction, blockHash []byte) ([]byte, error) {
	return t.Get(s.dirs[ACCOUNT.Name].Pack(tuple.Tuple{blockHash})).Get()
}

func (s *Storage) PutAccountUnclaimed(t fdb.Transaction, accountId, blockHash []byte) error {
	t.Set(s.dirs[ACCOUNT_UNCLAIMED.Name].Pack(tuple.Tuple{accountId}), blockHash)
	return nil
}

func (s *Storage) GetAccountUnclaimed(t fdb.Transaction, blockHash []byte) ([]byte, error) {
	return t.Get(s.dirs[ACCOUNT_UNCLAIMED.Name].Pack(tuple.Tuple{blockHash})).Get()
}

func (s *Storage) DeleteAccountUnclaimed(t fdb.Transaction, accountId []byte) error {
	t.Clear(s.dirs[ACCOUNT_UNCLAIMED.Name].Pack(tuple.Tuple{accountId}))
	return nil
}

func (s *Storage) PutBlock(t fdb.Transaction, blockHash, block []byte) error {
	t.Set(s.dirs[BLOCK.Name].Pack(tuple.Tuple{blockHash}), block)
	return nil
}

func (s *Storage) GetBlock(t fdb.Transaction, blockHash []byte) ([]byte, error) {
	return t.Get(s.dirs[BLOCK.Name].Pack(tuple.Tuple{blockHash})).Get()
}

func (s *Storage) PutUnclaimedBlock(t fdb.Transaction, blockHash, nextBlockHash []byte) error {
	t.Set(s.dirs[UNCLAIMED_BLOCK.Name].Pack(tuple.Tuple{blockHash}), nextBlockHash)
	return nil
}

func (s *Storage) GetUnclaimedBlock(t fdb.Transaction, blockHash []byte) ([]byte, error) {
	return t.Get(s.dirs[UNCLAIMED_BLOCK.Name].Pack(tuple.Tuple{blockHash})).Get()
}

func (s *Storage) DeleteUnclaimedBlock(t fdb.Transaction, blockHash []byte) error {
	t.Clear(s.dirs[UNCLAIMED_BLOCK.Name].Pack(tuple.Tuple{blockHash}))
	return nil
}

func (s *Storage) PutUnclaimedInfoBlock(t fdb.Transaction, blockHash []byte, unclaimedInfoBlock proto.Message) error {
	data, err := proto.Marshal(unclaimedInfoBlock)
	if err != nil {
		return err
	}
	t.Set(s.dirs[UNCLAIMED_INFO_BLOCK.Name].Pack(tuple.Tuple{blockHash}), data)
	return nil
}

func (s *Storage) DeleteUnclaimedInfoBlock(t fdb.Transaction, blockHash []byte) error {
	t.Clear(s.dirs[UNCLAIMED_INFO_BLOCK.Name].Pack(tuple.Tuple{blockHash}))
	return nil
}

func (s *Storage) PutClaimedBlock(t fdb.Transaction, blockHash []byte, time int64) error {
	t.Set(s.dirs[CLAIMED_BLOCK.Name].Pack(tuple.Tuple{blockHash}), []byte(fmt.Sprintf("%d", time)))
	return nil
}

func (s *Storage) PutNode(t fdb.Transaction, nodeId []byte, node proto.Message) error {
	data, err := proto.Marshal(node)
	if err != nil {
		return err
	}
	t.Set(s.dirs[NODE.Name].Pack(tuple.Tuple{nodeId}), data)
	return nil
}

func (s *Storage) PutNodeAddress(t fdb.Transaction, nodeId []byte, nodeAddress proto.Message) error {
	data, err := proto.Marshal(nodeAddress)
	if err != nil {
		return err
	}
	t.Set(s.dirs[NODE_ADDRESS.Name].Pack(tuple.Tuple{nodeId}), data)
	return nil
}
