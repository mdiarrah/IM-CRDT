package IpfsLink

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/pnet"
	"github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
)

const RENDEZ_VOUS = "MLS Test"
const TOPIC_NAME = "MLS Test Delivery Service"

const DHT_PROTOCOL_ID = "/mls-test/swarm/1"

func DHT_PKI_KEY(name string, id int) string {
	return "/keypkg/" + name + "/" + string(rune(id))
}

type Client struct {
	Ctx  context.Context
	Name string

	Host host.Host
	Dht  *dht.IpfsDHT

	Ps    *pubsub.PubSub
	Topic *pubsub.Topic
	Sub   *pubsub.Subscription

	Id string
}

/** Custom validator to allow anything to be stored in DHT (e.g. key packages) */
type CustomValidator struct{}

func (*CustomValidator) Validate(key string, value []byte) error {
	return nil
}
func (*CustomValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}

func InitClient(name string, bootstrapPeer string) *Client {
	ctx := context.Background()

	key := "edd99a84bbdd5c9cfc06bcc039d219b1000885ecba26901c02e7c8792bfaaa70"
	s := ""
	s += fmt.Sprintln("/key/swarm/psk/1.0.0/")
	s += fmt.Sprintln("/base16/")
	s += fmt.Sprintf("%s", key)
	psk, err := pnet.DecodeV1PSK(bytes.NewBuffer([]byte(s)))

	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(libp2p.EnableRelay(), libp2p.PrivateNetwork(psk))

	if err != nil {
		panic(fmt.Errorf("IPFSLink - InitClient, could not Create libP2P host\nerror: %s", err))
	}

	fmt.Println("Node ID:", host.ID())
	// fmt.Println("Addresses:")
	// for _, addr := range host.Addrs() {
	// 	fmt.Println("-", addr)
	// }
	fmt.Println("Full Address", host.Addrs()[0].String()+"/p2p/"+host.ID().String())
	idd := host.Addrs()[0].String() + "/p2p/" + host.ID().String()
	// If bootstrap node specified connect to it and then connect to peers
	//	using rendez-vous point
	// https://github.com/libp2p/go-libp2p/tree/master/examples/chat-with-rendezvous
	if bootstrapPeer != "" {
		fmt.Println("Bootstrap: Connecting to", bootstrapPeer, "...")

		peerInfo, err := peer.AddrInfoFromP2pAddr(multiaddr.StringCast(bootstrapPeer))
		if err != nil {
			panic(fmt.Errorf("IPFSLink - InitClient, could not retrieve Adrresses info from IFPS\nerror: %s\nBootstrapPeer : %s", err, bootstrapPeer))
		}

		if host.Connect(ctx, *peerInfo); err != nil {
			panic(fmt.Errorf("IPFSLink - InitClient, could not connect to another peer \nerror: %s", err))
		} else {
			fmt.Println("Bootstrap done !")
		}
	}

	// var dhtMode dht.Option
	// if bootstrapPeer == "" {
	// 	dhtMode = dht.Mode(dht.ModeServer)
	// } else {
	// 	dhtMode = dht.Mode(dht.ModeServer)
	// }
	dhtMode := dht.Mode(dht.ModeServer)

	dhtV, err := dht.New(ctx, host, dhtMode,
		dht.ProtocolPrefix(protocol.ID(DHT_PROTOCOL_ID)),
		dht.RoutingTableRefreshPeriod(5*time.Second))

	if err != nil {
		panic(fmt.Errorf("IPFSLink - InitClient, could not create the DHT \nerror: %s", err))
	}
	dhtV.Validator = &CustomValidator{}

	if err := dhtV.Bootstrap(ctx); err != nil {
		panic(fmt.Errorf("IPFSLink - InitClient, could not connect to the bootstrap\nerror: %s", err))
	}

	ps, topic, sub := SetupPubSub(host, ctx, TOPIC_NAME)
	P2PsetupDiscovery(host, ctx, dhtV)

	return &Client{
		Ctx:   ctx,
		Name:  name,
		Host:  host,
		Dht:   dhtV,
		Ps:    ps,
		Topic: topic,
		Sub:   sub,
		Id:    idd,
	}
}
func (client *Client) Close() {
	// if err := client.topic.Close(); err != nil {
	// 	panic(fmt.Errorf("IPFSLink - Close, could not Close Topic\nerror: %s", err))
	// }

	if err := client.Host.Close(); err != nil {
		panic(fmt.Errorf("IPFSLink - Close, could not Close Host\nerror: %s", err))
	}
}

func (client *Client) Publish(message []byte) {
	fmt.Println("Publising message", string(message))
	if err := client.Topic.Publish(client.Ctx, message); err != nil {
		fmt.Println("Error Sending Message:", err)
	}
}

/** Register at rendez-vous point and connect to found peers */
func P2PsetupDiscovery(host host.Host, ctx context.Context, dht *dht.IpfsDHT) {
	routingDiscovery := routing.NewRoutingDiscovery(dht)
	routingDiscovery.Advertise(ctx, RENDEZ_VOUS)
	fmt.Println("Advertising started")

	peerChan, err := routingDiscovery.FindPeers(ctx, RENDEZ_VOUS)
	if err != nil {
		panic(fmt.Errorf("Issue Starting RouteDiscovery to Findpeer :%s", err))
	}

	go func() {
		for peer := range peerChan {
			if peer.ID == host.ID() {
				continue
			}
			if err := host.Connect(ctx, peer); err != nil {
				fmt.Println("Could not connect to peer", err)
			} else {
				fmt.Println("Connected to peer", peer.ID, "(", peer.Addrs[0], "...)")
			}
		}
	}()
	fmt.Println("Advertising is running now")
}

/** Setup the PubSub and subscribe to appropriate topic */
func SetupPubSub(host host.Host, ctx context.Context, top string) (*pubsub.PubSub, *pubsub.Topic, *pubsub.Subscription) {
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(fmt.Errorf("IPFSLink - SetUpdPubSub, could not Create the GossipSub \nerror: %s", err))
	}

	topic, err := ps.Join(top)
	if err != nil {
		panic(fmt.Errorf("IPFSLink - SetUpdPubSub, could not Join the topic :%s \nerror: %s", TOPIC_NAME, err))
	}

	sub, err := topic.Subscribe()
	if err != nil {
		panic(fmt.Errorf("IPFSLink - SetUpdPubSub, could not Subscribe to topic \nerror: %s", err))
	}

	return ps, topic, sub
}
