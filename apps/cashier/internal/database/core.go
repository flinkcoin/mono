package database

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/flinkcoin/mono/libs/shared/pkg/base"
	"sync"
)

type ColumnFamily struct {
	ID   int
	Name string
}

const (
	ACCOUNT              = ColumnFamily{1, "account"}
	ACCOUNT_UNCLAIMED    = ColumnFamily{2, "account_unclaimed"}
	BLOCK                = ColumnFamily{3, "block"}
	CLAIMED_BLOCK        = ColumnFamily{4, "claimed_block"}
	UNCLAIMED_INFO_BLOCK = ColumnFamily{5, "unclaimed_info_block"}
	UNCLAIMED_BLOCK      = ColumnFamily{6, "unclaimed_block"}
	WEIGHT               = ColumnFamily{7, "weight"}
	NODE                 = ColumnFamily{8, "node"}
	NODE_ADDRESS         = ColumnFamily{9, "node_address"}
)

type Core struct {
	db   fdb.Database
	dirs map[string]directory.DirectorySubspace
	mu   sync.Mutex
}

func NewCore() (*Core, error) {
	fdb.MustAPIVersion(620)
	db := fdb.MustOpenDefault()

	dirs := make(map[string]directory.DirectorySubspace)
	root, err := directory.CreateOrOpen(db, []string{"mydb"}, nil)
	if err != nil {
		return nil, err
	}

	for _, cf := range []ColumnFamily{ACCOUNT, ACCOUNT_UNCLAIMED, BLOCK, CLAIMED_BLOCK, UNCLAIMED_INFO_BLOCK, UNCLAIMED_BLOCK, WEIGHT, NODE, NODE_ADDRESS} {
		dir, err := root.CreateOrOpen(db, []string{cf.Name}, nil)
		if err != nil {
			return nil, err
		}
		dirs[cf.Name] = dir
	}

	return &Core{
		db:   db,
		dirs: dirs,
	}, nil
}

func (s *Core) Put(cf ColumnFamily, key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(s.dirs[cf.Name].Pack(tuple.Tuple{key}), value)
		return nil, nil
	})
	return err
}

func (s *Core) Get(cf ColumnFamily, key []byte) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, err := s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		return tr.Get(s.dirs[cf.Name].Pack(tuple.Tuple{key})).Get()
	})
	if err != nil {
		return nil, err
	}
	return result.([]byte), nil
}

func (s *Core) Close() {
	// No explicit close needed for FoundationDB
}

func (s *Core) Open() {
	err := fdb.APIVersion(710)
	if err != nil {
		base.Log.Error("Failed to set API version: %v", "Error", err)
	}
	db, err := fdb.OpenDefault()

	if err != nil {
		base.Log.Error("Error opening database.", "err", err)
	}
	s.db = &db

	base.Log.Info("Successfully connected to FoundationDB")
}

// Write stores a key-value pair in the database
func (s *Core) Write(key string, value []byte) error {
	// Ensure the database is open
	if s.db == nil {
		return fmt.Errorf("database not initialized; call Open() first")
	}

	// Run a transaction to write the key-value pair
	_, err := s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(fdb.Key(key), value)
		return nil, nil
	})
	if err != nil {
		base.Log.Error("Failed to write key %s: %v", key, err)
		return err
	}

	base.Log.Info("Successfully wrote key!", "key", key)
	return nil
}

// Read retrieves the value for a given key from the database
func (s *Core) Read(key string) ([]byte, error) {
	// Ensure the database is open
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized; call Open() first")
	}

	// Run a transaction to read the value
	result, err := s.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		future := tr.Get(fdb.Key(key))
		return future.MustGet(), nil
	})
	if err != nil {
		base.Log.Error("Failed to read key %s: %v", key, err)
		return nil, err
	}

	// Check if the result is nil (key not found)
	if result == nil {
		base.Log.Info("Key %s not found", "key", key)
		return nil, nil
	}

	value, ok := result.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected type for value of key %s", key)
	}

	base.Log.Info("Successfully read key", "key", key)
	return value, nil
}
