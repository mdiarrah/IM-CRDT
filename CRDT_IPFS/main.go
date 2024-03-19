package main

import (
	"IPFS_CRDT/CRDTDag"
	"IPFS_CRDT/Payload"
	Tests "IPFS_CRDT/example/tests"
	IPFSLink "IPFS_CRDT/ipfsLink"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	// "time"

	files "github.com/ipfs/go-ipfs-files"
	// "github.com/pkg/profile"
)

// func testIPFSLink() {

// 	fmt.Println("-- Getting an IPFS node running -- ")

// 	ipfs1, err := IPFSLink.InitNode("peeerName ", "")

// 	fmt.Println("-- have an IPFS node running , now adding a file-- ")

// 	if err != nil {
// 		panic(fmt.Errorf("failed to intantiate an IPFS Daemon: %s", err))
// 	}

// 	ipfsPath, err := IPFSLink.AddIPFS(ipfs1, []byte("i'm connected to IPFS yahooo!!!\n"))

// 	fmt.Println("-- added an IPFS file now create another node -- ")

// 	if err != nil {
// 		panic(fmt.Errorf("failed to add file to IPFS %s", err))
// 	}
// 	time.Sleep(20 * time.Second)
// 	ipfs2, err := IPFSLink.InitNode("peerName2", ipfs1.Cr.Host.Addrs()[0].String()+"/p2p/"+ipfs1.Cr.Host.ID().String())
// 	fmt.Println("-- added a node, now get a file -- ")

// 	if err != nil {
// 		panic(fmt.Errorf("failed to intantiate a second IPFS Daemon: %s", err))
// 	}
// 	fmt.Println("encoded CID :", string(ipfsPath.Cid().Bytes()), "vs CID", ipfsPath.Cid().String())
// 	fil, err := IPFSLink.GetIPFS(ipfs2, ipfsPath.Cid().Bytes())

// 	go func() {
// 		for {
// 			msg, err := ipfs1.Cr.Sub.Next(ipfs1.Cr.Ctx)
// 			if err != nil {
// 				fmt.Println("\x1b[31m"+"Pub sub returned error, Aborting:", err, "\x1b[0m")
// 				break
// 			} else if msg.ReceivedFrom != ipfs1.Cr.Host.ID() {
// 				fmt.Println("Received message from", msg.ReceivedFrom,
// 					"data:", string(msg.Data))

// 				fmt.Println()
// 			} else {
// 				fmt.Println("Received message from myself :", msg.ReceivedFrom,
// 					"data:", string(msg.Data))
// 			}
// 		}
// 	}()

// 	go func() {
// 		for {
// 			msg, err := ipfs2.Cr.Sub.Next(ipfs2.Cr.Ctx)
// 			if err != nil {
// 				fmt.Println("\x1b[31m"+"Pub sub returned error, Aborting:", err, "\x1b[0m")
// 				break
// 			} else if msg.ReceivedFrom != ipfs2.Cr.Host.ID() {
// 				fmt.Println("Received message from", msg.ReceivedFrom,
// 					"data:", string(msg.Data))

// 				fmt.Println()
// 			} else {
// 				fmt.Println("Received message from myself :", msg.ReceivedFrom,
// 					"data:", string(msg.Data))
// 			}
// 		}
// 	}()

// 	fmt.Println("-- done-- ")

// 	err = files.WriteTo(fil, "./file.data")

// 	if err != nil {
// 		panic(fmt.Errorf("could not write out the fetched CID: %s", err))
// 	}
// 	//
// 	// "github.com/libp2p/go-libp2p-core/peer"

// 	fmt.Println("out - size : ", len(ipfs1.Topics))
// 	fmt.Println("out - size2 : ", len(ipfs2.Topics))
// 	fmt.Println("list peers1 : ", ipfs1.Cr.Ps.ListPeers(ipfs1.Cr.Sub.Topic())) //.ListPeers("bonjour"))
// 	fmt.Println("list peers2 : ", ipfs2.Cr.Ps.ListPeers(ipfs2.Cr.Sub.Topic()))
// 	//.ListPeers("bonjour"))

// 	time.Sleep(2 * time.Second)
// 	IPFSLink.PubIPFS(ipfs1, []byte("I'm connected"))
// 	time.Sleep(2 * time.Second)
// 	IPFSLink.PubIPFS(ipfs2, []byte("I'm connected too"))
// 	time.Sleep(20 * time.Second)
// }

// func testIPFSLink2() {

// 	fmt.Println("-- Getting an IPFS node running -- ")

// 	ipfs1, err := IPFSLink.InitNode("peeerName ", "")

// 	fmt.Println("-- have an IPFS node running , now adding a file-- ")

// 	if err != nil {
// 		panic(fmt.Errorf("failed to intantiate an IPFS Daemon: %s", err))
// 	}

// 	ipfsPath, err := IPFSLink.AddIPFS(ipfs1, []byte("i'm connected to IPFS yahooo!!!\n"))

// 	fmt.Println("-- added an IPFS file now create another node -- ")

// 	if err != nil {
// 		panic(fmt.Errorf("failed to add file to IPFS %s", err))
// 	}
// 	time.Sleep(20 * time.Second)
// 	bootstrapPeer := ipfs1.Cr.Host.Addrs()[0].String() + "/p2p/" + ipfs1.Cr.Host.ID().String()
// 	ipfs2, err := IPFSLink.InitNode("peerName2", bootstrapPeer)
// 	fmt.Println("-- added a node, now get a file -- ")
// 	man1 := CRDTDag.Create_CRDTManager(ipfs1, ".", "")
// 	man2 := CRDTDag.Create_CRDTManager(ipfs2, ".", bootstrapPeer)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to intantiate a second IPFS Daemon: %s", err))
// 	}
// 	fil, err := man2.GetNodeFromEncodedCid(man1.EncodeCid(ipfsPath))
// 	files.WriteTo(fil, "file1.data")
// }

type PayloadExample struct {
	X int
}

func (this *PayloadExample) FromString(payload string) {
	v, err := strconv.Atoi(payload)
	this.X = v
	if err != nil {
		panic(fmt.Errorf("error with PayloadExample FromString: %s", payload))
	}
	if false {

		files.WriteTo(nil, "file1.data")
		IPFSLink.InitClient("a", "a")

	}
}
func (this *PayloadExample) ToString() string {
	return strconv.Itoa(this.X)
}
func testCRDTDag() {
	x := CRDTDag.CRDTDagNode{}
	dd := make([]CRDTDag.EncodedStr, 0)
	dd = append(dd, CRDTDag.EncodedStr{Str: []byte("dependence1")})
	dd = append(dd, CRDTDag.EncodedStr{Str: []byte("dependence2")})
	peerID := "123PID321"
	pl := PayloadExample{X: 3}
	var plprime Payload.Payload = &pl
	x.CreateNode(dd, peerID, &plprime)
	x.ToFile("thisisafile.txt")
	y := CRDTDag.CRDTDagNode{}
	pl2 := PayloadExample{X: 3}
	var plprime2 Payload.Payload = &pl2
	y.CreateNodeFromFile("thisisafile.txt", &plprime2)
	y.ToFile("thisisafile2.txt")
}

var mu sync.Mutex

func main() {

	// Tests.RemoteTestSet()n

	peerName := flag.String("name", "FIRST", "name/identity of the current node")
	mode := flag.String("mode", "", "mode of the current application")
	updatesNB := flag.Int("updatesNB", 1000, "Number of updates")
	updating := flag.Bool("updating", false, "do I update the data")
	measurement := flag.Bool("TimeMeasurement", true, "do I Measure the different timefor each CID, adds files continainning these")
	ntpServ := flag.String("NTPS", "0.europe.pool.ntp.org", "Available NTP server for time measures")
	encode := flag.String("encode", "", "Data encription key")
	bootstrapPeer := flag.String("ni", "", "Client bootstrap for pubsub")
	IPFSbootstrap := flag.String("IPFSBootstrap", "", "IPFS bootstrap peer to have a private network")

	_ = measurement
	_ = updating

	flag.Parse()
	if *mode == "BootStrap" {
		fmt.Println("bootstrap peer :", *bootstrapPeer)
		// if err := os.Mkdir(*peerName, os.ModePerm); err != nil {
		// 	panic(err)
		// }
		if err := os.Mkdir(*peerName, os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}
		if err := os.Mkdir(*peerName+"/remote", os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}
		if err := os.Mkdir(*peerName+"/rootNode", os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}
		if err := os.Mkdir(*peerName+"/time", os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}
		Tests.Peer1Concu(*peerName, *updatesNB, *ntpServ, *encode, *measurement) // ------------- MANAGE CONCURENCY !!!
		// Tests.Peer1IPFS(*peerName, *updatesNB, *ntpServ) // ------------- NO CONCURENCY, ONLY IPFS ALONE !!!
		// Tests.Peer1(*peerName, *updatesNB, *ntpServ) // ------------- NO CONCURENCY, CRDT + IPFS  !!!
	} else if *mode == "update" {

		fmt.Println("bootstrap peer :", *bootstrapPeer)
		// if err := os.Mkdir(*peerName, os.ModePerm); err != nil {
		// 	panic(err)
		// }
		if err := os.Mkdir(*peerName+"/remote", os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}
		if err := os.Mkdir(*peerName+"/rootNode", os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}
		if err := os.Mkdir(*peerName+"/time", os.ModePerm); err != nil {
			fmt.Print(err, "\n")
		}

		// defer profile.Start(profile.CPUProfile).Stop()
		// if false {
		// 	Tests.Peer2Concu(*peerName, *bootstrapPeer, *updatesNB)
		// 	fmt.Println("test2")
		// }

		if *updating {
			fmt.Println("UPDATING IN FACT")
			Tests.Peer2ConcuUpdate(*peerName, *bootstrapPeer, *IPFSbootstrap, *updatesNB, *ntpServ, *encode, *measurement) // ------------- MANAGE CONCURENCY !!!
		} else {
			fmt.Println("NOT UPDATING FIOU")
			Tests.Peer2Concu(*peerName, *bootstrapPeer, *IPFSbootstrap, *updatesNB, *ntpServ, *encode, *measurement) // ------------- MANAGE CONCURENCY !!!
		}

		// Tests.Peer2IPFS(*peerName, *bootstrapPeer, *updatesNB, *ntpServ) // ------------- NO CONCURENCY, ONLY IPFS ALONE !!!
		// Tests.Peer2(*peerName, *bootstrapPeer, *updatesNB, *ntpServ) // ------------- NO CONCURENCY, CRDT + IPFS  !!!
	}
}
