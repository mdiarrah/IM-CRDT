package tests

// import (
// 	"IPFS_CRDT/example/Counter"
// 	IpfsLink "IPFS_CRDT/ipfsLink"
// 	"fmt"
// 	"time"
// )

// func TestCounterBootstrap() {
// 	sys1, err := IpfsLink.InitNode("BOOTSTRAP", "")
// 	if err != nil {
// 		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
// 	}
// 	CounterCrdt1 := Counter.Create_CRDTCounterOpBasedDag(sys1, "./node1", "", "")
// 	bootstrapPeer := sys1.Cr.Host.Addrs()[0].String() + "/p2p/" + sys1.Cr.Host.ID().String()
// 	fmt.Println("BootStrapPeer:", bootstrapPeer)
// 	for {
// 		for i := 0; i < 3; i++ {
// 			time.Sleep(time.Minute)
// 		}
// 		CounterCrdt1.CheckUpdate()
// 		x := CounterCrdt1.Lookup()
// 		fmt.Println("New Value of the Counter:", x.Lookup())
// 	}
// }
