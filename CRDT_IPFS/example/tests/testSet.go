package tests

import (
	"IPFS_CRDT/example/Set"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"fmt"
	"time"
)

func TestSet() {
	sys, err := IpfsLink.InitNode("joe", "")
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt := Set.Create_CRDTSetOpBasedDag(sys, ".", "")

	(&SetCrdt).Add("x")
	fmt.Println("Value :")
	l := SetCrdt.Lookup()
	fmt.Println(l.Lookup())

	(&SetCrdt).Add("y")
	fmt.Println("Value :")
	l = SetCrdt.Lookup()
	fmt.Println(l.Lookup())

	(&SetCrdt).Add("z")
	fmt.Println("Value :")
	l = SetCrdt.Lookup()
	fmt.Println(l.Lookup())

	fmt.Println("decrementing")
	(&SetCrdt).Remove("x")
	fmt.Println("Value :")
	l = SetCrdt.Lookup()
	fmt.Println(l.Lookup())
	fmt.Println("the end !! :)")

}

func RemoteTestSet() {

	time_start := time.Now()
	sys1, err := IpfsLink.InitNode("FIRST", "")
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, "./node1", "")
	bootstrapPeer := sys1.Cr.Host.Addrs()[0].String() + "/p2p/" + sys1.Cr.Host.ID().String()
	sys2, err := IpfsLink.InitNode("SECOND", bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt2 := Set.Create_CRDTSetOpBasedDag(sys2, "./node2", bootstrapPeer)
	sys3, err := IpfsLink.InitNode("THIRD", bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt3 := Set.Create_CRDTSetOpBasedDag(sys3, "./node3", bootstrapPeer)

	duration := time.Since(time_start)
	fmt.Println("duration to create nodes :", duration.Seconds())

	fmt.Println("let's take a nap")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("what a wonderfull day")

	time_start = time.Now()
	time_start = time.Now()
	(&SetCrdt1).Add("x1")

	(&SetCrdt2).Add("x2")

	(&SetCrdt3).Add("x3")

	duration = time.Since(time_start)
	fmt.Println("duration to send 3 nodes updates :", duration.Seconds())

	fmt.Println("i'm tired")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("I slept well")

	time_start = time.Now()

	SetCrdt1.CheckUpdate()
	SetCrdt2.CheckUpdate()
	SetCrdt3.CheckUpdate()

	duration = time.Since(time_start)
	fmt.Println("duration to check 2 new updates per nodes :", duration.Seconds())

	time_start = time.Now()

	l := SetCrdt1.Lookup()
	fmt.Println("Value 1 :", l.Lookup())

	l = SetCrdt2.Lookup()
	fmt.Println("\nValue 2 :", l.Lookup())

	l = SetCrdt3.Lookup()
	fmt.Println("\nValue 3 :", l.Lookup())

	duration = time.Since(time_start)
	fmt.Println("duration to do 3 lookups :", duration.Seconds())

	sys4, err := IpfsLink.InitNode("FORTH", sys1.Cr.Host.Addrs()[0].String()+"/p2p/"+sys1.Cr.Host.ID().String())
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt4 := Set.Create_CRDTSetOpBasedDag(sys4, "./node4", bootstrapPeer)
	fmt.Println("i'm tired again")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("I slept better")

	(&SetCrdt1).Add("y1")

	fmt.Println("i'm tired again")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	fmt.Println("I slept better")

	SetCrdt1.CheckUpdate()
	SetCrdt2.CheckUpdate()
	SetCrdt3.CheckUpdate()
	SetCrdt4.CheckUpdate()

	l = SetCrdt1.Lookup()
	fmt.Println("Value 1 :", l.Lookup())

	l = SetCrdt2.Lookup()
	fmt.Println("\nValue 2 :", l.Lookup())

	l = SetCrdt3.Lookup()
	fmt.Println("\nValue 3 :", l.Lookup())

	l = SetCrdt4.Lookup()
	fmt.Println("\nValue 4 :", l.Lookup())

}
