package IpfsLink

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	libp2pIFPS "github.com/ipfs/kubo/core/node/libp2p"
	"github.com/ipfs/kubo/plugin/loader"
	"github.com/ipfs/kubo/repo/fsrepo"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "pubsub-chat-example"

// printErr is like fmt.Printf, but writes to stderr.
func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}

// defaultNick generates a nickname based on the $USER environment variable and
// the last 8 chars of a peer ID.
func defaultNick(p peer.ID) string {
	return fmt.Sprintf("%s-%s", os.Getenv("USER"), shortID(p))
}

// shortID returns the last 8 chars of a base58-encoded peer id.
func shortID(p peer.ID) string {
	pretty := p.Pretty()
	return pretty[len(pretty)-8:]
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.Addrs[0])
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

type IpfsLink struct {
	Cancel    context.CancelFunc
	Ctx       context.Context
	IpfsCore  icore.CoreAPI
	IpfsNode  *core.IpfsNode
	Topics    []*pubsub.Topic
	Hst       host.Host
	GossipSub *pubsub.PubSub
	Cr        *Client
}

func InitNode(peerName string, bootstrapPeer string) (*IpfsLink, error) {
	ct, cancl := context.WithCancel(context.Background())

	// Spawn a local peer using a temporary path, for testing purposes
	ipfsA, nodeA, err := spawnEphemeral(ct)
	if err != nil {
		panic(fmt.Errorf("failed to spawn peer node: %s", err))
	}

	h := InitClient(peerName, bootstrapPeer)
	ipfs := IpfsLink{
		Cancel:    cancl,
		Ctx:       ct,
		IpfsCore:  ipfsA,
		IpfsNode:  nodeA,
		Hst:       nodeA.PeerHost,
		GossipSub: h.Ps,
		Cr:        h,
	}

	//fmt.Println(ipfs.IpfsNode.Peerstore.PeerInfo(ipfs.IpfsNode.PeerHost.ID()))
	return &ipfs, err
}

var loadPluginsOnce sync.Once

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

var flagExp = flag.Bool("experimental", false, "enable experimental features")

func createTempRepo() (string, error) {
	repoPath, err := os.MkdirTemp("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(io.Discard, 2048)
	if err != nil {
		return "", err
	}

	// When creating the repository, you can define custom settings on the repository, such as enabling experimental
	// features (See experimental-features.md) or customizing the gateway endpoint.
	// To do such things, you should modify the variable `cfg`. For example:
	if *flagExp {
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#ipfs-filestore
		cfg.Experimental.FilestoreEnabled = true
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#ipfs-urlstore
		cfg.Experimental.UrlstoreEnabled = true
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#ipfs-p2p
		cfg.Experimental.Libp2pStreamMounting = true
		// https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md#p2p-http-proxy
		cfg.Experimental.P2pHttpProxy = true
		// See also: https://github.com/ipfs/kubo/blob/master/docs/config.md
		// And: https://github.com/ipfs/kubo/blob/master/docs/experimental-features.md
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

/// ------ Spawning the node

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2pIFPS.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}
	return core.NewNode(ctx, nodeOptions)
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context) (icore.CoreAPI, *core.IpfsNode, error) {
	var onceErr error
	loadPluginsOnce.Do(func() {
		onceErr = setupPlugins("")
	})
	if onceErr != nil {
		return nil, nil, onceErr
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	node, err := createNode(ctx, repoPath)
	if err != nil {
		return nil, nil, err
	}

	api, err := coreapi.NewCoreAPI(node)

	return api, node, err
}

func AddIPFS(ipfs *IpfsLink, message []byte) (icorepath.Resolved, error) {

	peerCidFile, err := ipfs.IpfsCore.Unixfs().Add(ipfs.Ctx,
		files.NewBytesFile(message))
	if err != nil {
		panic(fmt.Errorf("could not add File: %s", err))
	}
	go ipfs.IpfsCore.Dht().Provide(ipfs.Ctx, peerCidFile)
	// if err != nil {
	// 	panic(fmt.Errorf("Could not provide File - %s", err))
	// }
	return peerCidFile, err
}

type CID struct{ str string }

func GetIPFS(ipfs *IpfsLink, c []byte) (files.Node, error) {
	// str_CID, err := ContentIdentifier.Decode(c)
	_, str_CID, err := cid.CidFromBytes(c)
	if err != nil {
		panic(fmt.Errorf("could not conver string of CID : \"%s\"\nerror:%s", c, err))
	}
	fil, err := ipfs.IpfsCore.Unixfs().Get(ipfs.Ctx, icorepath.IpfsPath(str_CID))
	if err != nil {
		panic(fmt.Errorf("could not get file with CID - %s : %s", c, err))
	}

	return fil, err
}

func PubIPFS(ipfs *IpfsLink, msg []byte) {
	// msgBytes, err := json.Marshal(msg)
	// if err != nil {
	// 	return err
	// }

	// x := 0
	// for x = 0; x < len(ipfs.Topics); x++ {
	// 	if ipfs.Topics[x].String() == topic {
	// 		break
	// 	}
	// }
	// if x == len(ipfs.Topics) {
	// 	fmt.Println("size : ", len(ipfs.Topics))
	// 	fmt.Println("Topic : ", topic, " is not ", ipfs.Topics[0].String())
	// 	SubIPFS(ipfs, topic)
	// // }
	// err = ipfs.GossipSub.Publish(topic, []byte(msgBytes))

	// if err != nil {
	// 	fmt.Println("error with publish : ", err)
	// }
	ipfs.Cr.Publish(msg)
}
