package db

import (
	"context"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"testing"
)

func readerToString(reader io.Reader) (string, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read from reader: %w", err)
	}
	return string(bytes), nil
}

//func ExecCommandInContainer(ctx context.Context, container testcontainers.Container, cmd []string) (string, error) {
//	code, read, err := container.Exec(ctx, cmd)
//	if err != nil {
//		return "", fmt.Errorf("failed to execute command: %w", err)
//	}
//
//	if code != 0 {
//		return "", fmt.Errorf("command exited with code %d: %s", code, read)
//	}
//
//	return read, nil
//}

func TestFoundationDB(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "foundationdb/foundationdb:7.3.57",
		ExposedPorts: []string{"4500/tcp"},
		Env: map[string]string{
			"FDB_NETWORKING_MODE": "host", // Set the desired networking mode
		},
		WaitingFor: wait.ForLog("FDBD joined cluster"),
	}

	fdbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	// Get the logs of the container
	reader, err := fdbContainer.Logs(ctx)
	if err != nil {
		t.Fatalf("Error reading logs: %v", err)
	}
	defer reader.Close()

	// Print the logs
	logBytes, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Error reading log contents: %v", err)
	}
	fmt.Println("Container Logs:\n", string(logBytes))

	if err != nil {
		t.Fatal(err)
	}
	defer fdbContainer.Terminate(ctx)

	// Example usage of ExecCommandInContainer
	//output, err := ExecCommandInContainer(ctx, fdbContainer, []string{"fdbcli", "--exec", "status"})
	//if err != nil {
	//	t.Fatalf("Error executing command in container: %v", err)
	//}
	//fmt.Println("Command Output:\n", output)

	// Get the mapped port
	host, err := fdbContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mappedPort, err := fdbContainer.MappedPort(ctx, "4500")
	if err != nil {
		t.Fatal(err)
	}

	// Here you would typically connect to FoundationDB using the host and mapped port
	// For example:
	// fdb := foundationdb.Open(host, mappedPort.Port())
	// Perform your tests here...
	fmt.Println("FoundationDB container is running at", host, ":", mappedPort.Port())

	fdb.MustAPIVersion(730)

	// Open the default database from the system cluster
	connectionString := fmt.Sprintf("docker:docker@%s:%s", host, mappedPort.Port())
	db, errr := fdb.OpenWithConnectionString(connectionString)

	if errr != nil {
		log.Fatalf("Unable to connect to FDB (%v)", errr)
	}

	// Database reads and writes happen inside transactions
	ret, err := db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tr.Set(fdb.Key("hello"), []byte("world"))
		return tr.Get(fdb.Key("foo")).MustGet(), nil
		// db.Transact automatically commits (and if necessary,
		// retries) the transaction
	})
	if err != nil {
		log.Fatalf("Unable to perform FDB transaction (%v)", err)
	}

	fmt.Printf("hello is now world, foo was: %s\n", string(ret.([]byte)))
}
