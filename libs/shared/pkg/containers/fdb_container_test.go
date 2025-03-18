package containers

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/stretchr/testify/assert"
)

func TestFdbContainer_CreatesClusterFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fdbContainer, err := NewFdbContainer(ctx, "")
	if err != nil {
		t.Fatalf("Failed to create FoundationDB container: %v", err)
	}

	clusterFile, err := fdbContainer.GetClusterFilePath(ctx)
	if err != nil {
		t.Fatalf("Failed to get cluster file path: %v", err)
	}
	assert.FileExists(t, clusterFile, "Cluster file should exist after container start")

	// Stop the container and verify the cluster file is deleted
	if err := fdbContainer.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
	_, err = os.Stat(clusterFile)
	assert.True(t, os.IsNotExist(err), "Cluster file should be deleted after container stop")
}

func TestFdbContainer_ShouldExecuteTransactions(t *testing.T) {
	err := fdb.APIVersion(710)
	if err != nil {
		t.Fatalf("Failed to set FDB API version: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fdbContainer, err := NewFdbContainer(ctx, "")
	if err != nil {
		t.Fatalf("Failed to create FoundationDB container: %v", err)
	}
	defer func() {
		if err := fdbContainer.Stop(ctx); err != nil {
			t.Errorf("Failed to stop container: %v", err)
		}
	}()

	clusterFile, err := fdbContainer.GetClusterFilePath(ctx)
	if err != nil {
		t.Fatalf("Failed to get cluster file path: %v", err)
	}

	db, err := fdb.OpenDatabase(clusterFile)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Test getting a non-existent key
	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		val, err := tr.Get(fdb.Key([]byte("bla"))).Get()
		assert.NoError(t, err)
		assert.Nil(t, val, "Value for non-existent key 'bla' should be nil")
		return nil, nil
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	// Set and get a value
	key := tuple.Tuple{"hello"}.Pack()
	value := tuple.Tuple{"world"}.Pack()
	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(fdb.Key(key), value)
		return nil, nil
	})
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	retrieved, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		return tr.Get(fdb.Key(key)).Get()
	})
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	assert.Equal(t, value, retrieved, "Retrieved value should match set value")
}

func TestFdbContainer_ShouldWorkWithReuse(t *testing.T) {
	err := fdb.APIVersion(710)
	if err != nil {
		t.Fatalf("Failed to set FDB API version: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Note: Testcontainers-go does not natively support reuse like Testcontainers for Java.
	// Simulate reuse by creating a new container instance and ensuring it works.
	fdbContainer, err := NewFdbContainer(ctx, "")
	if err != nil {
		t.Fatalf("Failed to create FoundationDB container: %v", err)
	}
	defer func() {
		if err := fdbContainer.Stop(ctx); err != nil {
			t.Errorf("Failed to stop container: %v", err)
		}
	}()

	clusterFile, err := fdbContainer.GetClusterFilePath(ctx)
	if err != nil {
		t.Fatalf("Failed to get cluster file path: %v", err)
	}

	db, err := fdb.OpenDatabase(clusterFile)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(fdb.Key(tuple.Tuple{"hello"}.Pack()), tuple.Tuple{"world"}.Pack())
		return nil, nil
	})
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}
}

func TestFdbContainer_ShouldRunWithSpecificVersion(t *testing.T) {
	err := fdb.APIVersion(710)
	if err != nil {
		t.Fatalf("Failed to set FDB API version: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	image := "foundationdb/foundationdb:7.1.61"
	fdbContainer, err := NewFdbContainer(ctx, image)
	if err != nil {
		t.Fatalf("Failed to create FoundationDB container: %v", err)
	}
	defer func() {
		if err := fdbContainer.Stop(ctx); err != nil {
			t.Errorf("Failed to stop container: %v", err)
		}
	}()

	clusterFile, err := fdbContainer.GetClusterFilePath(ctx)
	if err != nil {
		t.Fatalf("Failed to get cluster file path: %v", err)
	}

	db, err := fdb.OpenDatabase(clusterFile)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	assert.NotNil(t, db, "Database should be opened successfully")
}

func TestFdbContainer_Example(t *testing.T) {
	err := fdb.APIVersion(710)
	if err != nil {
		t.Fatalf("Failed to set FDB API version: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fdbContainer, err := NewFdbContainer(ctx, "")
	if err != nil {
		t.Fatalf("Failed to create FoundationDB container: %v", err)
	}
	defer func() {
		if err := fdbContainer.Stop(ctx); err != nil {
			t.Errorf("Failed to stop container: %v", err)
		}
	}()

	clusterFile, err := fdbContainer.GetClusterFilePath(ctx)
	if err != nil {
		t.Fatalf("Failed to get cluster file path: %v", err)
	}

	db, err := fdb.OpenDatabase(clusterFile)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	_, err = db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(fdb.Key(tuple.Tuple{"hello"}.Pack()), tuple.Tuple{"world"}.Pack())
		return nil, nil
	})
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}
}
