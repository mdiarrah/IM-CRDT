package tests

import (
	"IPFS_CRDT/example/Counter"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"fmt"
	"time"
)

func TestCounter() {
	sys, err := IpfsLink.InitNode("joe", "")
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	CounterCrdt := Counter.Create_CRDTCounterOpBasedDag(sys, ".", "")

	(&CounterCrdt).Increment()
	fmt.Println("Value :")
	l := CounterCrdt.Lookup()
	fmt.Println(l.Lookup())

	(&CounterCrdt).Increment()
	fmt.Println("Value :")
	l = CounterCrdt.Lookup()
	fmt.Println(l.Lookup())

	(&CounterCrdt).Increment()
	fmt.Println("Value :")
	l = CounterCrdt.Lookup()
	fmt.Println(l.Lookup())

	fmt.Println("decrementing")
	(&CounterCrdt).Decrement()
	fmt.Println("Value :")
	l = CounterCrdt.Lookup()
	fmt.Println(l.Lookup())
	fmt.Println("the end !! :)")

}

func RemoteTestCounter() {

	time_start := time.Now()
	sys1, err := IpfsLink.InitNode("FIRST", "")
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	CounterCrdt1 := Counter.Create_CRDTCounterOpBasedDag(sys1, "./node1", "")
	bootstrapPeer := sys1.Cr.Host.Addrs()[0].String() + "/p2p/" + sys1.Cr.Host.ID().String()
	sys2, err := IpfsLink.InitNode("SECOND", bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	CounterCrdt2 := Counter.Create_CRDTCounterOpBasedDag(sys2, "./node2", bootstrapPeer)
	sys3, err := IpfsLink.InitNode("THIRD", bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	CounterCrdt3 := Counter.Create_CRDTCounterOpBasedDag(sys3, "./node3", bootstrapPeer)

	duration := time.Since(time_start)
	fmt.Println("duration to create nodes :", duration.Seconds())

	fmt.Println("let's take a nap")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("what a wonderfull day")

	time_start = time.Now()
	time_start = time.Now()
	(&CounterCrdt1).Increment()

	(&CounterCrdt2).Increment()

	(&CounterCrdt3).Decrement()
	duration = time.Since(time_start)
	fmt.Println("duration to send 3 nodes updates :", duration.Seconds())

	fmt.Println("i'm tired")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("I slept well")

	time_start = time.Now()

	CounterCrdt1.CheckUpdate()
	CounterCrdt2.CheckUpdate()
	CounterCrdt3.CheckUpdate()

	duration = time.Since(time_start)
	fmt.Println("duration to check 2 new updates per nodes :", duration.Seconds())

	time_start = time.Now()

	l := CounterCrdt1.Lookup()
	fmt.Println("Value 1 :", l.Lookup())

	l = CounterCrdt2.Lookup()
	fmt.Println("\nValue 2 :", l.Lookup())

	l = CounterCrdt3.Lookup()
	fmt.Println("\nValue 3 :", l.Lookup())

	duration = time.Since(time_start)
	fmt.Println("duration to do 3 lookups :", duration.Seconds())

	sys4, err := IpfsLink.InitNode("FORTH", sys1.Cr.Host.Addrs()[0].String()+"/p2p/"+sys1.Cr.Host.ID().String())
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	CounterCrdt4 := Counter.Create_CRDTCounterOpBasedDag(sys4, "./node4", bootstrapPeer)
	fmt.Println("i'm tired again")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("I slept better")

	(&CounterCrdt1).Increment()

	fmt.Println("i'm tired again")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("I slept better")

	CounterCrdt1.CheckUpdate()
	CounterCrdt2.CheckUpdate()
	CounterCrdt3.CheckUpdate()
	CounterCrdt4.CheckUpdate()

	l = CounterCrdt1.Lookup()
	fmt.Println("Value 1 :", l.Lookup())

	l = CounterCrdt2.Lookup()
	fmt.Println("\nValue 2 :", l.Lookup())

	l = CounterCrdt3.Lookup()
	fmt.Println("\nValue 3 :", l.Lookup())

	l = CounterCrdt4.Lookup()
	fmt.Println("\nValue 4 :", l.Lookup())

}
