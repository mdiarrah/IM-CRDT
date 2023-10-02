package tests

import (
	"IPFS_CRDT/example/Set"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func WriteFile(fileName string, b []byte) {
	fil, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Errorf("WRITEFILE - could Not Open RootNode to update rootnodefolder\nerror: %s", err))
	}
	_, err = fil.Write(b)
	if err != nil {
		panic(fmt.Errorf("could Not write in RootNode to WRITEFILE - \nerror: %s", err))
	}
	err = fil.Close()
	if err != nil {
		panic(fmt.Errorf("could Not Close - WRITEFILE\nerror: %s", err))
	}
}

func Peer1(peername string, nbUpdates int) {
	sys1, err := IpfsLink.InitNode("BOOTSTRAP", "")
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, "")
	str := ""
	for i := range sys1.Cr.Host.Addrs() {
		s := sys1.Cr.Host.Addrs()[i].String()
		str += s + "/p2p/" + sys1.Cr.Host.ID().String() + "\n"
	}
	if _, err := os.Stat("./ID2"); !errors.Is(err, os.ErrNotExist) {
		os.Remove("./ID2")
	}
	WriteFile("./ID2", []byte(str))
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
	}
	file, err := os.OpenFile(peername+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}
	ti := time.Now()
	k := 0
	for k < nbUpdates {
		time.Sleep(30 * time.Millisecond)

		strList := SetCrdt1.CheckUpdate()

		for j := 0; j < len(strList); j++ {
			file.WriteString(strList[j].Cid + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + "\n")
		}
		x := SetCrdt1.Lookup()
		fmt.Println("New Value of the Set:", x.Lookup())

		if time.Since(ti) > 500*time.Millisecond {
			t := strconv.FormatInt(time.Now().UnixMilli(), 10)
			encodedCid := SetCrdt1.Add(sys1.Cr.Id + "VALUE ADDED" + strconv.Itoa(k))
			file.WriteString(encodedCid + "," + t + "\n")
			k++
			ti = time.Now()
		}
	}

	if err = file.Close(); err != nil {
		panic(fmt.Errorf("Error closing file\nerror : %s", err))
	}

	time.Sleep(5 * time.Minute)
}

func Peer2(folder_storage string, bootstrapPeer string, nbUpdates int) {
	sys1, err := IpfsLink.InitNode(folder_storage, bootstrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, folder_storage, bootstrapPeer)

	file, err := os.OpenFile(folder_storage+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}
	for {
		time.Sleep(30 * time.Millisecond)

		strList := SetCrdt1.CheckUpdate()

		for j := 0; j < len(strList); j++ {
			file.WriteString(strList[j].Cid + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + "\n")
		}

		x := SetCrdt1.Lookup()

		fmt.Println("New Value of the Set:", x.Lookup())
	}
}
