package tests

import (
	"IPFS_CRDT/example/Set"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
	// "github.com/beevik/ntp"
)

func GetTime(ntpServ string) int {
	// // x := "ssh root@" + ntpServ + " sh -c 'date +%s%N | cut -b1-13' "
	// // cmd := exec.Command(x, "./")
	// // out, err := cmd.Output()
	// // if err != nil {
	// // 	// if there was any error, print it here
	// // 	fmt.Println("could not run command: ", err)
	// // }

	// // ti, err := strconv.Atoi(string(out))
	// ti, err := ntp.Time(ntpServ)
	// if err != nil {
	// 	// n := 0
	// 	panic(fmt.Errorf("Failed To get time : %s", err))
	// 	// ti = time.Time{0}
	// 	// err = nil
	// 	// for n < 100 && err != nil { // try over 1 s to get ntp time
	// 	// 	time.Sleep(10 * time.Millisecond)
	// 	// 	ti, err = ntp.Time("2.fr.pool.ntp.org")
	// 	// 	n += 1
	// 	// }
	// 	// if err != nil { // if ntp is impossible to
	// 	// 	ti = time.Now()

	// 	// 	f, _ := os.OpenFile(peername+"/error.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	// 	// 	f.WriteString("ERROR NTP !!!!")
	// 	// 	f.Close()
	// 	// }
	// }
	return int(time.Now().UnixMilli())
}

// \/ BOOTSTRAP PEER IS THIS ONE \/
func Peer1Concu(peername string, nbUpdates int, ntpServ string) {
	sys1, err := IpfsLink.InitNode(peername, "")
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
	time.Sleep(20 * time.Second)

	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, "")

	file, err := os.OpenFile(peername+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time,time_retrieve,time_compute,time_add_IPFS\n")

	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}
	fmt.Println("Starting the Set, sleeping 30s to wait others")

	ti := time.Now()
	// Sleep 30s before emiting updates to wait others
	for time.Since(ti) < 300*time.Second {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}
	}
	fmt.Printf("Starting the Set, updating %d times\n", nbUpdates)
	ti = time.Now()
	k := 0
	for k < nbUpdates {
		time.Sleep(30 * time.Microsecond)
		strList := SetCrdt1.CheckUpdate()
		t := strconv.Itoa(GetTime(ntpServ))
		if len(strList) > 0 {

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}
		if time.Since(ti) >= time.Millisecond*1000 {
			t := strconv.Itoa(GetTime(ntpServ))
			time_start := time.Now()
			encodedCid := SetCrdt1.Add(sys1.Cr.Id + "VALUE ADDED" + strconv.Itoa(k))
			file.WriteString(encodedCid + "," + t + "," + "0,0," + strconv.Itoa(int(time.Since(time_start).Nanoseconds())) + "\n")
			k++
			ti = time.Now()
		}

	}
	if err = file.Close(); err != nil {
		panic(fmt.Errorf("Error closing file\nerror : %s", err))
	}
	for {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}
		// x := SetCrdt1.Lookup()
		// fmt.Println("New Value of the Set:", x.Lookup())
	}
}

func Peer2Concu(peername string, bootStrapPeer string, nbUpdates int, ntpServ string) {
	sys1, err := IpfsLink.InitNode(peername, bootStrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}

	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, bootStrapPeer)
	file, err := os.OpenFile(peername+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time,time_retrieve,time_compute,time_add_IPFS\n")
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
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}

	}
}

func Peer2ConcuUpdate(peername string, bootStrapPeer string, nbUpdates int, ntpServ string) {
	sys1, err := IpfsLink.InitNode(peername, bootStrapPeer)
	if err != nil {
		panic(fmt.Errorf("Failed To instanciate IFPS & LibP2P clients : %s", err))
	}

	SetCrdt1 := Set.Create_CRDTSetOpBasedDag(sys1, peername, bootStrapPeer)
	file, err := os.OpenFile(peername+"/time.csv", os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString("CID,time,time_retrieve,time_compute,time_add_IPFS\n")
	if err != nil {
		panic(fmt.Errorf("Error openning file file\nerror : %s", err))
	}

	// Sleep 30s before emiting updates to wait others
	ti := time.Now()
	for time.Since(ti) < 300*time.Second {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}
	}

	fmt.Printf("Starting the Set, updating %d times\n", nbUpdates)
	ti = time.Now()
	k := 0
	for k < nbUpdates {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}

		// x := SetCrdt1.Lookup()
		// fmt.Println("New Value of the Set:", x.Lookup())

		if time.Since(ti) >= time.Millisecond*1000 {
			time_start := time.Now()
			encodedCid := SetCrdt1.Add(sys1.Cr.Id + "VALUE ADDED" + strconv.Itoa(k))
			file.WriteString(encodedCid + "," + strconv.Itoa(GetTime(ntpServ)) + ",0,0," + strconv.Itoa(int(time.Since(time_start).Nanoseconds())) + "\n")
			k++
			ti = time.Now()
		}

	}
	if err = file.Close(); err != nil {
		panic(fmt.Errorf("Error closing file\nerror : %s", err))
	}
	for {
		time.Sleep(30 * time.Microsecond)

		strList := SetCrdt1.CheckUpdate()
		if len(strList) > 0 {
			t := strconv.Itoa(GetTime(ntpServ))

			for j := 0; j < len(strList); j++ {
				file.WriteString(strList[j].Cid + "," + t + "," + strconv.Itoa(strList[j].IntegrityCheckTime) + "," + strconv.Itoa(strList[j].CalculTime) + ",0\n")
			}
		}

		// x := SetCrdt1.Lookup()
		// fmt.Println("New Value of the Set:", x.Lookup())
	}
}
