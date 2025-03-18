package networking

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/flinkcoin/mono/apps/broker/internal/router"
	"github.com/flinkcoin/mono/libs/shared/pkg/base"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"time"
)

type Host struct {
	host   *host.Host
	router *router.Router
}

func NewHost(h *host.Host, r *Router) *Host {
	return &Host{
		host:   h,
		router: r,
	}
}

func (n *Host) Init() {
	// To construct a simple host with all the default settings, just use `New`

	// Now, normally you do not just want a simple host, you want
	// that is fully configured to best support your p2p application.
	// Let's create a second host setting some more options.

	// Set your own keypair
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}
	n.host, err = libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
		),
		// support TLS connections_
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),

		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}

	base.Log.Info("Hello World, my second hosts ID is %s\n", "hostKey:", n.host.ID())

	startListener(context.Background(), n.host)
}

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func (n *Host) startListener(ctx context.Context, ha host.Host) {
	fullAddr := getHostAddress(ha)
	log.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	ha.SetStreamHandler("/cashier/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if payload, err := readPayload(s); err != nil {
			log.Println(err)
			s.Reset()
		}

		peerId := s.Conn().RemotePeer()
		pubKey := s.Conn().RemotePublicKey(



		log.Printf("Received stream from peer: %s, public key: %s\n", peerId, pubKey)
		pubKey.Raw()
		n.router.ProcessMsg(payload, peerId, pubKey)
	})
}

func (n *Host) Connect(peerek string) {
	// Define the peer address to connect to
	peerAddr, err := ma.NewMultiaddr(peerek)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the peer ID from the multiaddress
	peerInfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the peer
	if err := n.host.Connect(context.Background(), *peerInfo); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected too", peerInfo.ID)

	// Create a new stream to the peer
	s, err := n.host.NewStream(context.Background(), peerInfo.ID, "/echo/1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		log.Fatal(err)
	}
}

func readPayload(s network.Stream) ([]byte, error) {
	// Create a buffered reader for the stream
	buf := bufio.NewReader(s)

	// Read the 4-byte length prefix (big-endian)
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(buf, lengthBuf)
	if err != nil {
		return nil, err // Could be EOF, network error, etc.
	}

	// Convert the 4 bytes to an integer
	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, errors.New("empty message length")
	}

	// Read exactly 'length' bytes for the message
	payload := make([]byte, length)
	_, err = io.ReadFull(buf, payload)
	if err != nil {
		return nil, err // Could be EOF or partial read
	}

	// Return the message bytes
	return payload, nil
}
