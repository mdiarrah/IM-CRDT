package tests

import (
	"IPFS_CRDT/example/NoConcurrency"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func Peer1IPFS(peername string, nbUpdates int, ntpServ string) {
	sys1, err := IpfsLink.InitNode(peername, "")
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt1 := NoConcurrency.InitIPFSSet(sys1, peername, "")
	str := ""
	for i := range sys1.Cr.Host.Addrs() {
		s := sys1.Cr.Host.Addrs()[i].String()
		str += s + "/p2p/" + sys1.Cr.Host.ID().String() + "\n"

	}
	if _, err := os.Stat("./ID2"); !errors.Is(err, os.ErrNotExist) {
		os.Remove("./ID2")
	}
	WriteFile("./ID2", []byte(str))
	time.Sleep(10 * time.Second)
	file, err := os.OpenFile(peername+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time,timeSend,timeRetrieve\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}

	ti := time.Now()
	time.Sleep(60 * time.Second)
	k := 0
	for k < nbUpdates {
		time.Sleep(30 * time.Microsecond)

		time_start := time.Now()
		strList := SetCrdt1.CheckUpdate()
		for j := 0; j < len(strList); j++ {
			file.WriteString(strList[j] + "," + strconv.Itoa(GetTime(ntpServ)) + ",0," + strconv.Itoa(int(time.Since(time_start).Nanoseconds())) + "\n")
		}
		// x := SetCrdt1.Lookup()
		// fmt.Println("New Value of the Set:", x)

		if time.Since(ti) >= time.Millisecond*1000 {
			t := strconv.Itoa(GetTime(ntpServ))
			time_start = time.Now()
			encodedCid := SetCrdt1.Add(sys1.Cr.Id + "VALUE ADDED" + strconv.Itoa(k))
			file.WriteString(string(encodedCid) + "," + t + "," + strconv.Itoa(int(time.Since(time_start).Nanoseconds())) + ",0\n")
			k++
			ti = time.Now()
		}

	}

	if err = file.Close(); err != nil {
		panic(fmt.Errorf("Error closing file\nerror : %s", err))
	}

	time.Sleep(5 * time.Minute)

}

func Peer2IPFS(folder_storage string, bootstrapPeer string, nbUpdates int, ntpServ string) {
	sys1, err := IpfsLink.InitNode(folder_storage, bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt1 := NoConcurrency.InitIPFSSet(sys1, folder_storage, bootstrapPeer)

	file, err := os.OpenFile(folder_storage+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time,timeSend,timeRetrieve\n")

	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}
	for {
		time.Sleep(30 * time.Microsecond)
		time_start := time.Now()
		strList := SetCrdt1.CheckUpdate()
		for j := 0; j < len(strList); j++ {
			file.WriteString(strList[j] + "," + strconv.Itoa(GetTime(ntpServ)) + ",0," + strconv.Itoa(int(time.Since(time_start).Nanoseconds())) + "\n")
		}

		// x := SetCrdt1.Lookup()

		// fmt.Println("New Value of the Set:", x)
	}
}
