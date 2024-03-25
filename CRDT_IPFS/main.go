package main

import (
	"IPFS_CRDT/CRDTDag"
	"IPFS_CRDT/Payload"
	IPFSLink "IPFS_CRDT/ipfsLink"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/ipfs/go-log/v2"

	// "time"

	files "github.com/ipfs/go-ipfs-files"
	// "github.com/pkg/profile"
)

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

var logger = log.Logger("IPFS_CRDT")

func main() {
	log.SetAllLoggers(log.LevelDebug)
	log.SetLogLevel("IPFS_CRDT", "debug")
	isBootstrap := flag.Bool("bootstrap", false, "flag to indicate whether its is a bootstrap node")
	bootstrapPath := flag.String("path", "", "Bootstrap peer json file")
	flag.Parse()

	bootstrapConfig := IPFSLink.MultiAddressesJson{
		AddressList: make([]string, 0),
	}

	if !*isBootstrap {
		if _, err := os.Stat(*bootstrapPath); err == nil {
			jsonData, err := ioutil.ReadFile(*bootstrapPath)
			if err != nil {
				logger.Errorf("error loading bootstrap file, %v", err)
				return
			}
			err = json.Unmarshal(jsonData, &bootstrapConfig)
			if err != nil {
				logger.Errorf("error Unmarshalling bootstrap file, %v", err)
				return
			}
		} else {
			logger.Errorf("error, bootstrap file '%s' file doesn't exist %v", *bootstrapPath, err)
			return
		}
	}

	fmt.Println("bootstrap peer :", bootstrapConfig.AddressList)
	if _, err := IPFSLink.CreateRepoAndConfig(bootstrapConfig.AddressList); err != nil {
		panic(err)
	}

	if err := IPFSLink.CreateAndStartIpfsNode(bootstrapConfig.AddressList); err != nil {
		panic(err)
	}

	select {}

}
