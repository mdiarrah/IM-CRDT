package tests

import (
	"IPFS_CRDT/example/Set"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/semaphore"
	// "github.com/beevik/ntp"
)

func GetTime(ntpServ string) int {
	return int(time.Now().UnixMilli())
}

func getSema(sema *semaphore.Weighted, ctx context.Context) {
	t := time.Now()
	err := sema.Acquire(ctx, 1)
	for err != nil && time.Since(t) < 10*time.Second {
		time.Sleep(10 * time.Microsecond)
		err = sema.Acquire(ctx, 1)
	}
	if err != nil {
		panic(fmt.Errorf("Semaphore of READ/WRITE file locked !!!!\n Cannot acquire it\n"))
	}
}

func returnSema(sema *semaphore.Weighted) {
	sema.Release(1)
}

// \/ BOOTSTRAP PEER IS THIS ONE \/
func Peer1Concu(peername string, nbUpdates int, ntpServ string, encode string, measurement bool) {

	fileRead, err := os.OpenFile(peername+"/time/FileRead.log", os.O_CREATE|os.O_WRONLY, 0755)
	file, err := os.OpenFile(peername+"/time/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	sema := semaphore.NewWeighted(1)

	sys1, err := IpfsLink.InitNode(peername, "", make([]byte, 0))
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}

	str := ""
	for i := range sys1.Cr.Host.Addrs() {
		s := sys1.Cr.Host.Addrs()[i].String()
		str += s + "/p2p/" + sys1.Cr.Host.ID().String() + "\n"
	}
	if _, err := os.Stat("./ID2"); !errors.Is(err, os.ErrNotExist) {
		os.Remove("./ID2")
	}

	WriteFile("./ID2", []byte(str))

	bytesIPFS_Node, err := sys1.IpfsNode.Peerstore.PeerInfo(sys1.IpfsNode.Identity).MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("Failed To Marshall IFPS Identity & LibP2P clients : %s", err))
	}
	if _, err := os.Stat("./IDBootstrapIPFS"); !errors.Is(err, os.ErrNotExist) {
		os.Remove("./IDBootstrapIPFS")
	}
	WriteFile("./IDBootstrapIPFS", bytesIPFS_Node)

	time.Sleep(20 * time.Second)

	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, "", encode, measurement)

	fileRead.WriteString("Taking Sema to write headers ... ")
	getSema(sema, sys1.Ctx)
	file.WriteString("CID,time,time_retrieve,time_compute,time_add_IPFS,time_encrypt,time_decrypt,time_Retreive_Whole_Batch\n")
	returnSema(sema)
	fileRead.WriteString("Header just written\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}
	fmt.Println("Starting the Set, sleeping 30s to wait others")

	ti := time.Now()
	// Sleep 300s before emiting updates to wait others
	for time.Since(ti) < 300*time.Second {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			fileRead.WriteString("Just Received some updates\n")
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				getSema(sema, sys1.Ctx)
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].RetrievalAlone) + "," + strconv.Itoa(strList[j].CalculTime) + ",0,0," + strconv.Itoa(strList[j].Time_decrypt) + "," + strconv.Itoa(strList[j].RetrievalTotal) + "\n")
				returnSema(sema)
				fileRead.WriteString("writing 1 line\n")
			}
			fileRead.WriteString("all update received are handled\n= = = = = = =\n")
		}
	}

	fmt.Printf("Starting the Set, updating %d times\n", nbUpdates)
	ti = time.Now()

	// Send updates concurrently every 1 seconds
	go sendUpdates(nbUpdates, &SetCrdt1, ntpServ, file, sys1.Cr.Id, sema)

	//regularly scan files if there is any new received updates
	for {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			fileRead.WriteString("Just Received some updates\n")
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				getSema(sema, sys1.Ctx)
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].RetrievalAlone) + "," + strconv.Itoa(strList[j].CalculTime) + ",0,0," + strconv.Itoa(strList[j].Time_decrypt) + "," + strconv.Itoa(strList[j].RetrievalTotal) + "\n")
				returnSema(sema)
				fileRead.WriteString("writing 1 line\n")
			}
			fileRead.WriteString("all update received are handled\n= = = = = = =\n")
		}
		// x := SetCrdt1.Lookup()
		// fmt.Println("New Value of the Set:", x.Lookup())
	}
}

func Peer2Concu(peername string, bootStrapPeer string, IPFSbootStrapPeer string, nbUpdates int, ntpServ string, encode string, measurement bool) {
	IPFSbootstrapBytes, err := os.ReadFile(IPFSbootStrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To Read IFPS bootstrap peer multiaddr : %s", err))
	}
	sys1, err := IpfsLink.InitNode(peername, bootStrapPeer, IPFSbootstrapBytes)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	time.Sleep(10 * time.Second)

	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, bootStrapPeer, encode, measurement)
	file, err := os.OpenFile(peername+"/time/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time,time_retrieve,time_compute,time_add_IPFS,time_encrypt,time_decrypt,time_Retreive_Whole_Batch\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}
	fmt.Printf("Starting the Set, updating %d times\n", nbUpdates)
	for {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].RetrievalAlone) + "," + strconv.Itoa(strList[j].CalculTime) + ",0,0," + strconv.Itoa(strList[j].Time_decrypt) + "," + strconv.Itoa(strList[j].RetrievalTotal) + "\n")
			}
		}

	}
}

func Peer2ConcuUpdate(peername string, bootStrapPeer string, IPFSbootStrapPeer string, nbUpdates int, ntpServ string, encode string, measurement bool) {
	sema := semaphore.NewWeighted(1)
	IPFSbootstrapBytes, err := os.ReadFile(IPFSbootStrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To Read IFPS bootstrap peer multiaddr : %s", err))
	}
	sys1, err := IpfsLink.InitNode(peername, bootStrapPeer, IPFSbootstrapBytes)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}
	time.Sleep(10 * time.Second)
	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, bootStrapPeer, encode, measurement)
	file, err := os.OpenFile(peername+"/time/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	fileRead, err := os.OpenFile(peername+"/time/FileRead.log", os.O_CREATE|os.O_WRONLY, 0755)
	fileRead.WriteString("Taking Sema to write headers ... ")
	getSema(sema, sys1.Ctx)
	file.WriteString("CID,time,time_retrieve,time_compute,time_add_IPFS,time_encrypt,time_decrypt,time_Retreive_Whole_Batch\n")
	returnSema(sema)
	fileRead.WriteString("Header just written\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}

	// Sleep 300s before emiting updates to wait others
	ti := time.Now()
	for time.Since(ti) < 300*time.Second {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			fileRead.WriteString("Just Received some updates\n")
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				getSema(sema, sys1.Ctx)
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].RetrievalAlone) + "," + strconv.Itoa(strList[j].CalculTime) + ",0,0," + strconv.Itoa(strList[j].Time_decrypt) + "," + strconv.Itoa(strList[j].RetrievalTotal) + "\n")
				returnSema(sema)
				fileRead.WriteString("writing 1 line\n")
			}
			fileRead.WriteString("all update received are handled\n= = = = = = =\n")
		}
	}

	// Send updates concurrently every 1 seconds
	go sendUpdates(nbUpdates, &SetCrdt1, ntpServ, file, sys1.Cr.Id, sema)

	//regularly scan files if there is any new received updates
	fmt.Printf("Starting the Set, updating %d times\n", nbUpdates)
	ti = time.Now()
	k := 0
	for k < nbUpdates {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			fileRead.WriteString("Just Received some updates\n")
			t := strconv.Itoa(GetTime(ntpServ))
			for j := 0; j < len(strList); j++ {
				getSema(sema, sys1.Ctx)
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].RetrievalAlone) + "," + strconv.Itoa(strList[j].CalculTime) + ",0,0," + strconv.Itoa(strList[j].Time_decrypt) + "," + strconv.Itoa(strList[j].RetrievalTotal) + "\n")
				returnSema(sema)
				fileRead.WriteString("writing 1 line\n")
			}
			fileRead.WriteString("all update received are handled\n= = = = = = =\n")
		}
	}
	if err := file.Close(); err != nil {
		panic(fmt.Errorf("Error closing file\nerror : %s", err))
	}

}

func sendUpdates(nbUpdates int, SetCrdt1 *Set.CRDTSetOpBasedDag, ntpServ string, file *os.File, netID string, sema *semaphore.Weighted) {
	fileWrite, _ := os.OpenFile(SetCrdt1.GetCRDTManager().Nodes_storage_enplacement+"/time/FileWrite.log", os.O_CREATE|os.O_WRONLY, 0755)
	fileWrite.WriteString(fmt.Sprintf("Starting the Set, updating %d times\n", nbUpdates))
	ti := time.Now()
	k := 0
	for k < nbUpdates {
		time.Sleep(30 * time.Microsecond)

		if time.Since(ti) >= time.Millisecond*1000 {
			fileWrite.WriteString("updating the data\n")
			encodedCid, times := SetCrdt1.Add(netID + "VALUE ADDED" + strconv.Itoa(k))
			fileWrite.WriteString("updating the data - taking sema\n")
			getSema(sema, context.Background())
			fileWrite.WriteString("Semaphore tooken\n")
			file.WriteString(encodedCid + "," + strconv.Itoa(GetTime(ntpServ)) + "," + "0,0," + strconv.Itoa(times.Time_add) + "," + strconv.Itoa(times.Time_encrypt) + ",0,0\n")
			fileWrite.WriteString("returning Semaphore\n")
			returnSema(sema)
			fileWrite.WriteString("WRITE - 1 line added to time.csv\n")
			k++
			ti = time.Now()
		}

	}
	fileWrite.WriteString("WRITE - all updates are done\n")
	fileWrite.Close()
}
