package IpfsLink

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"

	//formatIPFS "github.com/ipfs/go-ipld-format"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/bootstrap"
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
	Cancel          context.CancelFunc
	Ctx             context.Context
	IpfsCore        icore.CoreAPI
	IpfsNode        *core.IpfsNode
	Topics          []*pubsub.Topic
	Hst             host.Host
	GossipSub       *pubsub.PubSub
	Cr              *Client
	paralelRetrieve bool
}

func InitNode(peerName string, bootstrapPeer string, ipfsBootstrap []byte) (*IpfsLink, error) {
	ct, cancl := context.WithCancel(context.Background())

	// Spawn a local peer using a temporary path, for testing purposes
	var idBootstrap peer.AddrInfo
	var ipfsA icore.CoreAPI
	var nodeA *core.IpfsNode
	var err error

	if len(ipfsBootstrap) > 0 {

		e := idBootstrap.UnmarshalJSON(ipfsBootstrap)
		if e != nil {
			panic(fmt.Errorf("couldn't Unmarshal bootstrap peer addr info, error : %s", e))
		}
		ipfsA, nodeA, err = spawnEphemeral(ct, &idBootstrap)
	} else {
		ipfsA, nodeA, err = spawnEphemeral(ct, nil)
	}

	if err != nil {
		panic(fmt.Errorf("failed to spawn peer node: %s", err))
	}
	h := InitClient(peerName, bootstrapPeer)
	ipfs := IpfsLink{
		Cancel:          cancl,
		Ctx:             ct,
		IpfsCore:        ipfsA,
		IpfsNode:        nodeA,
		Hst:             nodeA.PeerHost,
		GossipSub:       h.Ps,
		Cr:              h,
		paralelRetrieve: false,
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

func createTempRepo(btstrap []peer.AddrInfo) (string, error) {
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

	// cfg.Addresses.Gateway = config.Strings{"/ip4/172.16.192.10/tcp/8080"}
	// cfg.Addresses.API = config.Strings{"/ip4/172.16.192.10/tcp/5001"}

	cfg.SetBootstrapPeers(btstrap)

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

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}
	return node, nil

}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context, btstrap *peer.AddrInfo) (icore.CoreAPI, *core.IpfsNode, error) {
	var onceErr error
	loadPluginsOnce.Do(func() {
		onceErr = setupPlugins("")
	})
	if onceErr != nil {
		return nil, nil, onceErr
	}

	// Create a Temporary Repo
	var m []peer.AddrInfo

	if btstrap != nil {
		m = make([]peer.AddrInfo, 1)
		m[0] = *btstrap
	} else {
		m = make([]peer.AddrInfo, 0)
	}

	repoPath, err := createTempRepo(m)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	// Create an IPFS node
	printErr("repository : %s\n", repoPath)
	os.WriteFile(repoPath+"/swarm.key", []byte("/key/swarm/psk/1.0.0/\n/base16/\nedd99a84bbdd5c9cfc06bcc039d219b1000885ecba26901c02e7c8792bfaaa70"), fs.FileMode(os.O_CREATE|os.O_WRONLY|os.O_APPEND))

	node, err := createNode(ctx, repoPath)
	if err != nil {
		return nil, nil, err
	}

	node.PNetFingerprint = []byte("4c7dc2a2735a84b4b11ff5b39225aa771cea1abd3acf9b98708a25f286df851c")
	// Connect the node to the other private network nodes

	var bstcfg bootstrap.BootstrapConfig
	if btstrap != nil {

		bstcfg = bootstrap.BootstrapConfig{

			MinPeerThreshold:  1,
			Period:            60 * time.Second,
			ConnectionTimeout: 30 * time.Second,
			BootstrapPeers: func() []peer.AddrInfo {
				m := make([]peer.AddrInfo, 1)
				m[0] = *btstrap
				return m
			},
		}
	} else {
		bstcfg = bootstrap.BootstrapConfig{
			MinPeerThreshold:  0,
			Period:            60 * time.Second,
			ConnectionTimeout: 30 * time.Second,
			BootstrapPeers: func() []peer.AddrInfo {
				m := make([]peer.AddrInfo, 0)
				// m[0] = node.Peerstore.PeerInfo(node.Identity)
				return m
			},
		}
	}

	node.Bootstrap(bstcfg)

	api, err := coreapi.NewCoreAPI(node)

	if btstrap != nil {
		api.Swarm().Connect(ctx, *btstrap)
	}

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

func GetIPFS(ipfs *IpfsLink, cids [][]byte) ([]files.Node, error) {
	// str_CID, err := ContentIdentifier.Decode(c)
	var files []files.Node = make([]files.Node, len(cids))
	var StrCids []cid.Cid = make([]cid.Cid, len(cids))
	var err error
	var file *os.File

	ti := time.Now()
	if len(cids) > 0 {

		file, err = os.OpenFile("node1/time/timeConcurrentRetrieve.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(fmt.Errorf("could not Close Debug File in IPFSLink:: GetIPFS\nerror:%s", err))
		}

		file.WriteString("" +
			"===============================New Batch of Cid To Retrieve===============================\n" +
			"=============================The Cids to  Download are thoose=============================\n")
	}

	for index, c := range cids {
		_, cid, err := cid.CidFromBytes(c)
		if err != nil {
			panic(fmt.Errorf("could not conver string of CID : \"%s\"\nerror:%s", c, err))
		}

		StrCids[index] = cid
		file.WriteString(cid.String() + "\n")
	}

	wg := sync.WaitGroup{}
	wg.Add(len(cids))
	errhapened := true
	for errhapened {
		errhapened = false
		for index, c := range cids {
			clocal := c
			if ipfs.paralelRetrieve {
				go func(i int) {
					str_CID := StrCids[i]
					cctx, _ := context.WithDeadline(ipfs.Ctx, time.Now().Add(time.Second*30))
					//files[i], err = ipfs.IpfsCore.Dag().Get(cctx, str_CID)
					files[i], err = ipfs.IpfsCore.Unixfs().Get(cctx, icorepath.IpfsPath(str_CID))

					if err != nil {
						printErr("could not get file with CID - %s : %s", clocal, err)
						errhapened = true
					}
					wg.Done()
				}(index)
			} else {
				str_CID := StrCids[index]
				file.WriteString(fmt.Sprintf("Asking the CID %s \n", str_CID))
				cctx, _ := context.WithDeadline(ipfs.Ctx, time.Now().Add(time.Second*30))
				f, err2 := ipfs.IpfsNode.DAG.Get(cctx, str_CID)
				cctx, _ = context.WithDeadline(ipfs.Ctx, time.Now().Add(time.Second*30))
				files[index], err = ipfs.IpfsCore.Unixfs().Get(cctx, icorepath.IpfsPath(str_CID))
				if err != nil {
					printErr("could not get file with CID - %s : %s\n", clocal, err)
					printErr("what we got from IPFSNODE.dag  error :  %s\n  data : %s\n====================================\n", err2, f)
					errhapened = true
				}
			}
		}
		if ipfs.paralelRetrieve {
			wg.Wait()
		}
	}
	file.WriteString("Got all the cids asked\n")

	if len(cids) > 0 {
		file.WriteString("\n" +
			"Nb of Cids: " + strconv.Itoa(len(cids)) + "\n" +
			"Time To Download: " + strconv.FormatInt(time.Since(ti).Milliseconds(), 10) + " ms\n" +
			"=================================The end of CID retrieval=================================\n" +
			"\n" +
			"\n" +
			"\n")
		err = file.Close()
		if err != nil {
			panic(fmt.Errorf("could not Close Debug File in IPFSLink:: GetIPFS\nerror:%s", err))
		}
	} else {

		file.WriteString("\n" +
			fmt.Sprintf("Even if no CID Where downloaded, len(cids):%d", len(cids)) +
			"=================================The end of CID retrieval=================================\n")
	}

	return files, err
}

func PubIPFS(ipfs *IpfsLink, msg []byte) {
	ipfs.Cr.Publish(msg)
}
