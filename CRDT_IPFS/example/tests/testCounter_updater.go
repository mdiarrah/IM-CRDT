package tests

import (
	"IPFS_CRDT/example/Counter"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"fmt"
	"math/rand"
	"time"
)

func TestCounterUpdater(i string, bootstrapPeer string) {
	sys1, err := IpfsLink.InitNode(i, bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	CounterCrdt1 := Counter.Create_CRDTCounterOpBasedDag(sys1, i, bootstrapPeer)

	rand.Seed(int64(int(i[0]) * rand.Int()))

	incr := 10 + (rand.Int() % 20)
	decr := 10 + (rand.Int() % 20)
	fmt.Println("I am going to  increment", incr, "time and to decrement", decr, "times in a random order with random sleeptimes between 3 and 100 seconds")
	for {
		for i := 0; i < rand.Int()%60+3; i++ {
			time.Sleep(time.Second)
		}
		CounterCrdt1.CheckUpdate()
		if rand.Int()%2 == 0 && incr > 0 {
			CounterCrdt1.Increment()
			incr--
		} else if decr > 0 {
			CounterCrdt1.Decrement()
			decr--
		}
		x := CounterCrdt1.Lookup()
		fmt.Println("New Value of the Counter:", x.Lookup())
	}
}
