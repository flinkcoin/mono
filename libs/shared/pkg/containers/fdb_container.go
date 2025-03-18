package containers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Constants
const (
	defaultImage = "foundationdb/foundationdb:7.1.61"
	socatImage   = "alpine/socat:1.8.0.1"
	internalPort = 4500
)

// FdbContainer represents a FoundationDB container with a Socat proxy.
type FdbContainer struct {
	container       testcontainers.Container
	proxy           testcontainers.Container
	network         *testcontainers.DockerNetwork
	bindPort        int
	clusterFilePath string
	networkAlias    string
	cancelCtx       context.CancelFunc
	initialized     bool
}

func NewFdbContainer(ctx context.Context, image string) (*FdbContainer, error) {
	if image == "" {
		image = defaultImage
	}

	// Create a new network
	network, err := network.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	// Generate a unique network alias
	networkAlias := "fdb-" + strconv.FormatInt(time.Now().UnixNano(), 16)

	// Create and start the Socat proxy container
	proxyReq := testcontainers.ContainerRequest{
		Image:        socatImage,
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", internalPort)},
		Networks:     []string{network.Name},
		Cmd:          []string{"tcp-listen:" + strconv.Itoa(internalPort+1) + ",fork,reuseaddr", "tcp:" + networkAlias + ":" + strconv.Itoa(internalPort+1)},
	}
	proxy, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: proxyReq,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy container: %w", err)
	}

	// Get the mapped port for Socat's internal port (4500)
	mappedPort, err := proxy.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", internalPort)))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}
	bindPort, err := strconv.Atoi(mappedPort.Port())
	if err != nil {
		return nil, fmt.Errorf("failed to parse bind port: %w", err)
	}

	// Create and start the FoundationDB container
	fdbReq := testcontainers.ContainerRequest{
		Image:    image,
		Networks: []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {networkAlias},
		},
		Env: map[string]string{
			"FDB_NETWORKING_MODE": "host",
			"FDB_PORT":            strconv.Itoa(bindPort),
		},
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", internalPort)},
		WaitingFor:   wait.ForLog("FDBD joined cluster."),
		Name:         "testcontainers-fdb-" + strconv.FormatInt(time.Now().UnixNano(), 16),
	}
	fdbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: fdbReq,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create FoundationDB container: %w", err)
	}

	// Proxy the bind port using Socat
	proxyCtx, cancelCtx := context.WithCancel(ctx)
	go func() {
		_, _, err := proxy.Exec(proxyCtx, []string{"socat", "TCP-LISTEN:4500,fork,reuseaddr", "TCP:" + networkAlias + ":" + strconv.Itoa(bindPort)})

		if proxyCtx.Err() != nil {
			log.Printf("Proxy container stopped outside: %v", proxyCtx.Err())
			return
		} else if err != nil {
			log.Printf("Error executing socat command: %v", err)
		}
	}()

	// Initialize the container struct
	fdb := &FdbContainer{
		container:    fdbContainer,
		proxy:        proxy,
		network:      network,
		bindPort:     bindPort,
		networkAlias: networkAlias,
		cancelCtx:    cancelCtx,
		initialized:  false,
	}

	// Initialize the database
	if err := fdb.initDatabase(ctx); err != nil {
		return nil, err
	}

	return fdb, nil
}

// GetConnectionString returns the connection string for the FoundationDB cluster.
func (fdb *FdbContainer) GetConnectionString(ctx context.Context) (string, error) {
	host, err := fdb.proxy.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get host: %w", err)
	}
	mappedPort, err := fdb.proxy.MappedPort(ctx, nat.Port(fmt.Sprintf("%d/tcp", internalPort)))
	if err != nil {
		return "", fmt.Errorf("failed to get mapped port: %w", err)
	}
	return fmt.Sprintf("docker:docker@%s:%s", host, mappedPort.Port()), nil
}

// GetClusterFilePath returns the path to a temporary cluster file containing the connection string.
// The file is created on the first call and deleted when the container stops.
func (fdb *FdbContainer) GetClusterFilePath(ctx context.Context) (string, error) {
	if fdb.clusterFilePath != "" {
		return fdb.clusterFilePath, nil
	}

	tmpFile, err := os.CreateTemp("", "fdb_*.cluster")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	connectionString, err := fdb.GetConnectionString(ctx)
	if err != nil {
		return "", err
	}
	if _, err := tmpFile.WriteString(connectionString); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	fdb.clusterFilePath = tmpFile.Name()
	log.Printf("Using cluster file %s", fdb.clusterFilePath)
	return fdb.clusterFilePath, nil
}

// Stop terminates the containers and cleans up resources.
func (fdb *FdbContainer) Stop(ctx context.Context) error {
	log.Printf("Stopping FoundationDB container...")

	fdb.cancelCtx()

	if fdb.clusterFilePath != "" {
		if err := os.Remove(fdb.clusterFilePath); err != nil {
			log.Printf("Failed to remove cluster file: %v", err)
		}
	}
	if err := fdb.container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate FoundationDB container: %w", err)
	}
	if err := fdb.proxy.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate proxy container: %w", err)
	}
	if err := fdb.network.Remove(ctx); err != nil {
		return fmt.Errorf("failed to remove network: %w", err)
	}
	return nil
}

// initDatabase initializes the database if it is not already initialized.
func (fdb *FdbContainer) initDatabase(ctx context.Context) error {

	if err := fdb.initDatabaseSingleInMemory(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	initialized, err := fdb.isDatabaseInitialized(ctx)
	if err != nil {
		return fmt.Errorf("failed to check database status: %w", err)
	}
	if !initialized {
		if err := fdb.initDatabaseSingleInMemory(ctx); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
	}
	return nil
}

// isDatabaseInitialized checks if the database is already initialized.
//func (fdb *FdbContainer) isDatabaseInitialized(ctx context.Context) (bool, error) {
//	output, err := fdb.runCliExecOutput(ctx, "status minimal")
//	if err != nil {
//		return false, err
//	}
//	return strings.Contains(output, "The database is available"), nil
//}

func (fdb *FdbContainer) isDatabaseInitialized(ctx context.Context) (bool, error) {
	return fdb.initialized, nil
}

// initDatabaseSingleInMemory initializes a single in-memory database.
func (fdb *FdbContainer) initDatabaseSingleInMemory(ctx context.Context) error {
	log.Printf("Initializing a single in-memory database...")
	output, err := fdb.runCliExecOutput(ctx, "configure new single memory")
	if err != nil {
		return err
	}
	if !strings.Contains(output, "Database created") {
		return fmt.Errorf("database not created: %s", output)
	}

	fdb.initialized = true

	log.Printf("Initialized successfully with a single in-memory database.")
	return nil
}

// runCliExecOutput executes an fdbcli command and returns its output.
func (fdb *FdbContainer) runCliExecOutput(ctx context.Context, command string) (string, error) {
	exitCode, reader, err := fdb.container.Exec(ctx, []string{"/usr/bin/fdbcli", "--exec", command})

	if err != nil {
		return "", fmt.Errorf("failed to execute fdbcli: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, reader)

	if err != nil {
		return "", fmt.Errorf("failed to read fdbcli output: %w", err)
	}

	log.Printf("fdbcli output: %s", stdoutBuf.String())

	if exitCode != 0 {
		return "", fmt.Errorf("fdbcli exited with code %d: %s", exitCode, stderrBuf.String())
	}

	return stdoutBuf.String(), nil
}
