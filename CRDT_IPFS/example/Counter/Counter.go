package Counter

import (
	CRDTDag "IPFS_CRDT/CRDTDag"
	CRDT "IPFS_CRDT/Crdt"
	Payload "IPFS_CRDT/Payload"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	// files "github.com/ipfs/go-ipfs-files"
)

// =======================================================================================
// Payload - OpBased
// =======================================================================================

type Operation int

const (
	INCREMENT Operation = iota
	DECREMENT
)

func (self Operation) ToString() string {
	str := ""
	if self == INCREMENT {
		str += "increment"
	} else if self == DECREMENT {
		str += "decrement"
	}
	return str
}
func (op *Operation) op_from_string(s string) {
	fmt.Println("op_from_string - operation received :", s)
	if s == "increment" {
		*op = INCREMENT
	} else if s == "decrement" {
		*op = DECREMENT
	}
}

type PayloadOpBased struct {
	Op Operation
	Id string
}

func (self *PayloadOpBased) Create_PayloadOpBased(s string, o1 Operation) {

	self.Op = o1
	self.Id = s
}
func (self *PayloadOpBased) ToString() string {

	str := self.Op.ToString()
	str += " " + self.Id
	return str
}
func (self *PayloadOpBased) FromString(s string) {

	res := strings.Split(s, " ")
	self.Op.op_from_string(res[0])
	self.Id = res[1]
}

// =======================================================================================
// CRDTCounter OpBased
// =======================================================================================

type CRDTCounterOpBased struct {
	sys         *IpfsLink.IpfsLink
	added       map[string]int
	decremented map[string]int
}

func Create_CRDTCounterOpBased(s *IpfsLink.IpfsLink) CRDTCounterOpBased {
	return CRDTCounterOpBased{
		sys:         s,
		added:       make(map[string]int, 0),
		decremented: make(map[string]int, 0),
	}
}

func (self *CRDTCounterOpBased) increment() {

	self.added[self.sys.Hst.ID().Pretty()] += 1
}
func (self *CRDTCounterOpBased) receivedIncrement(peerID string) {

	self.added[peerID] += 1
}

func (self *CRDTCounterOpBased) decrement() {

	self.decremented[self.sys.Hst.ID().Pretty()] += 1
}
func (self *CRDTCounterOpBased) receivedDecrement(peerID string) {

	self.decremented[peerID] += 1
}

func (self *CRDTCounterOpBased) Lookup() int {

	i := 0
	for x := range self.added {
		fmt.Println("add", self.added[x])

		i += self.added[x]
	}
	for x := range self.decremented {
		fmt.Println("decrease", self.decremented[x])
		i -= self.decremented[x]
	}
	return i
}

func (self *CRDTCounterOpBased) ToFile(file string) {

	i := ""
	for x := range self.added {
		i += x + "," + strconv.Itoa(self.added[x]) + " "
	}
	i += "\n"
	for x := range self.decremented {
		i += x + "," + strconv.Itoa(self.decremented[x]) + " "
	}
	f, err := os.Create(file)
	if err != nil {
		panic(fmt.Errorf("CRDTDagNode - ToFile Could not Create the file %s\nError: %s", file, err))
	}
	f.WriteString(i)
	err = f.Close()
	if err != nil {
		panic(fmt.Errorf("CRDTDagNode - ToFile Could not Write to the file %s\nError: %s", file, err))
	}
}

// =======================================================================================
// CRDTCounterDagNode OpBased
// =======================================================================================

type CRDTCounterOpBasedDagNode struct {
	DagNode CRDTDag.CRDTDagNode
}

func (self *CRDTCounterOpBasedDagNode) FromFile(fil string) {
	var pl Payload.Payload = &PayloadOpBased{Op: INCREMENT}
	self.DagNode.CreateNodeFromFile(fil, &pl)
}

func (self *CRDTCounterOpBasedDagNode) GetDirect_dependency() []CRDTDag.EncodedStr {

	return self.DagNode.DirectDependency
}

func (self *CRDTCounterOpBasedDagNode) ToFile(file string) {

	self.DagNode.ToFile(file)
}
func (self *CRDTCounterOpBasedDagNode) GetEvent() *Payload.Payload {

	return self.DagNode.Event
}
func (self *CRDTCounterOpBasedDagNode) GetPiD() string {

	return self.DagNode.PID
}
func (self *CRDTCounterOpBasedDagNode) CreateEmptyNode() *CRDTDag.CRDTDagNodeInterface {
	n := CreateDagNode(INCREMENT, "")
	var node CRDTDag.CRDTDagNodeInterface = &n
	return &node
}
func CreateDagNode(o Operation, id string) CRDTCounterOpBasedDagNode {
	var pl Payload.Payload = &PayloadOpBased{Op: o, Id: id}
	slic := make([]CRDTDag.EncodedStr, 0)
	return CRDTCounterOpBasedDagNode{
		DagNode: CRDTDag.CRDTDagNode{
			Event:            &pl,
			PID:              id,
			DirectDependency: slic,
		},
	}
}

// =======================================================================================
// CRDTCounterDag OpBased
// =======================================================================================

type CRDTCounterOpBasedDag struct {
	dag CRDTDag.CRDTManager
}

func (self *CRDTCounterOpBasedDag) GetDag() *CRDTDag.CRDTManager {

	return &self.dag
}
func (self *CRDTCounterOpBasedDag) SendRemoteUpdates() {

	self.dag.SendRemoteUpdates()
}
func (self *CRDTCounterOpBasedDag) GetCRDTManager() *CRDTDag.CRDTManager {

	return &self.dag
}
func (self *CRDTCounterOpBasedDag) Merge(cids []CRDTDag.EncodedStr) []string {
	//TODO Manage concurrency
	for _, cid := range cids {
		find := false
		for x := range self.dag.GetAllNodes() {
			if string(self.dag.GetAllNodes()[x]) == string(cid.Str) {
				find = true
				break
			}
		}
		if !find {
			// TODO HERE !!
			// fils, err := self.dag.GetNodeFromEncodedCid(append(make([]CRDTDag.EncodedStr, 0), cid))
			// if err != nil {
			// 	panic(fmt.Errorf("could not retrieve the node %s , error :%s", cid.Str, err))
			// }
			fstr := self.dag.NextFileName()
			if _, err := os.Stat(fstr); !errors.Is(err, os.ErrNotExist) {
				os.Remove(fstr)
			}
			// files.WriteTo(fils[0], fstr)
			n := CreateDagNode(INCREMENT, "")
			n.FromFile(fstr)

			self.remoteAddNode(cid, n)
		}
	}
	return nil
}
func (self *CRDTCounterOpBasedDag) remoteAddNode(cID CRDTDag.EncodedStr, newnode CRDTCounterOpBasedDagNode) {
	var pl CRDTDag.CRDTDagNodeInterface = &newnode
	self.dag.RemoteAddNodeSuper(cID, &pl)
}

func (self *CRDTCounterOpBasedDag) Increment() {
	newNode := CreateDagNode(INCREMENT, self.GetSys().Hst.ID().Pretty())
	for dependency := range self.dag.Root_nodes {
		// fmt.Println("dep:", self.dag.Root_nodes[dependency].Str)
		newNode.DagNode.DirectDependency = append(newNode.DagNode.DirectDependency, self.dag.Root_nodes[dependency])
	}

	strFile := self.dag.NextFileName()
	if _, err := os.Stat(strFile); !errors.Is(err, os.ErrNotExist) {
		os.Remove(strFile)
	}
	newNode.ToFile(strFile)
	bytes, err := os.ReadFile(strFile)
	if err != nil {
		panic(fmt.Errorf("ERROR INCREMENT CRDTCounterOpBasedDag, could not read file\nerror: %s", err))
	}
	path, err := IpfsLink.AddIPFS(self.dag.Sys, bytes)
	if err != nil {
		panic(fmt.Errorf("CRDTCounterOpBasedDag Increment, could not add the file to IFPS\nerror: %s", err))
	}

	encodedCid := self.dag.EncodeCid(path)

	// _, c, _ := cid.CidFromBytes(encodedCid.Str)
	// fmt.Println("encodedCid Increment :", c.String())
	var pl CRDTDag.CRDTDagNodeInterface = &newNode

	self.dag.AddNode(encodedCid, &pl) // TODOCounterCrdt Complete Node interface

	self.SendRemoteUpdates()

}
func (self *CRDTCounterOpBasedDag) Decrement() {

	newNode := CreateDagNode(DECREMENT, self.GetSys().Hst.ID().Pretty())
	for dependency := range self.dag.Root_nodes {
		newNode.DagNode.DirectDependency = append(newNode.DagNode.DirectDependency, self.dag.Root_nodes[dependency])
	}

	strFile := self.dag.NextFileName()
	if _, err := os.Stat(strFile); !errors.Is(err, os.ErrNotExist) {
		os.Remove(strFile)
	}
	newNode.ToFile(strFile)
	bytes, err := os.ReadFile(strFile)
	if err != nil {
		panic(fmt.Errorf("ERROR INCREMENT CRDTCounterOpBasedDag, could not read file\nerror: %s", err))
	}
	path, err := IpfsLink.AddIPFS(self.dag.Sys, bytes)
	if err != nil {
		panic(fmt.Errorf("CRDTCounterOpBasedDag Decrement, could not add the file to IFPS\nerror: %s", err))
	}

	encodedCid := self.dag.EncodeCid(path)
	// _, c, _ := cid.CidFromBytes(encodedCid.Str)
	// fmt.Println("encodedCid Decrement :", c.String())
	var pl CRDTDag.CRDTDagNodeInterface = &newNode
	self.dag.AddNode(encodedCid, &pl)
	self.SendRemoteUpdates()
}

func Create_CRDTCounterOpBasedDag(sys *IpfsLink.IpfsLink, storage_emplacement string, bootStrapPeer string, key string, measurement bool) CRDTCounterOpBasedDag {
	man := CRDTDag.Create_CRDTManager(sys, storage_emplacement, bootStrapPeer, key, measurement)
	crdtCounter := CRDTCounterOpBasedDag{dag: man}

	var pl CRDTDag.CRDTDag = &crdtCounter

	CRDTDag.CheckForRemoteUpdates(&pl, sys.Cr.Sub, man.Sys.Ctx)
	return crdtCounter
}

func (self *CRDTCounterOpBasedDag) GetSys() *IpfsLink.IpfsLink {

	return self.dag.Sys
}

func (self *CRDTCounterOpBasedDag) Lookup_ToSpecifyType() *CRDT.CRDT {

	crdt := CRDTCounterOpBased{
		sys:         self.GetSys(),
		added:       make(map[string]int, 0),
		decremented: make(map[string]int, 0),
	}
	for x := range self.dag.GetAllNodes() {
		node := self.dag.GetAllNodesInterface()[x]
		if (*(*node).GetEvent()).(*PayloadOpBased).Op == INCREMENT {
			// fmt.Println("add")
			crdt.receivedIncrement((*node).GetPiD())
		} else {
			// fmt.Println("remove")
			crdt.receivedDecrement((*node).GetPiD())
		}
	}
	var pl CRDT.CRDT = &crdt
	return &pl
}
func (self *CRDTCounterOpBasedDag) Lookup() CRDTCounterOpBased {
	// crdt := self.Lookup_ToSpecifyType()
	// var pl CRDTDag.CRDTDag = &crdtCounter
	return *(*self.Lookup_ToSpecifyType()).(*CRDTCounterOpBased)
}

func (self *CRDTCounterOpBasedDag) CheckUpdate() {
	files, err := ioutil.ReadDir(self.GetDag().Nodes_storage_enplacement + "/remote")
	if err != nil {
		panic(fmt.Errorf("CheckUpdate - Checkupdate could not open folder\nerror: %s", err))
	}
	to_add := make([]CRDTDag.EncodedStr, 0)
	for _, file := range files {

		fil, err := os.OpenFile(self.GetDag().Nodes_storage_enplacement+"/remote/"+file.Name(), os.O_RDONLY, os.ModeAppend)
		if err != nil {
			panic(fmt.Errorf("error in checkupdate, Could not open the sub file\nError: %s", err))
		}
		stat, err := fil.Stat()
		if err != nil {
			panic(fmt.Errorf("error in checkupdate, Could not get stat the sub file\nError: %s", err))
		}
		bytesread := make([]byte, stat.Size())
		n, err := fil.Read(bytesread)
		if err != nil {
			panic(fmt.Errorf("error in checkupdate, Could not read the sub file\nError: %s", err))
		}
		if int64(n) != stat.Size() {
			panic(fmt.Errorf("error in checkupdate, Could not read entirely the sub file\nError: read %d byte unstead of %d", n, stat.Size()))
		}
		err = fil.Close()
		if err != nil {
			panic(fmt.Errorf("error in checkupdate, Could not close the sub file\nError: %s", err))
		}
		err = os.Remove(self.GetDag().Nodes_storage_enplacement + "/remote/" + file.Name())
		if err != nil || errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf("error in checkupdate, Could not remove the sub file\nError: %s", err))
		}
		to_add = append(to_add, CRDTDag.EncodedStr{Str: bytesread})
	}
	self.Merge(to_add)

	self.GetDag().UpdateRootNodeFolder()
}

// Next function must be useless
func (self *CRDTCounterOpBasedDag) CheckRootNodes() {
	files, err := ioutil.ReadDir(self.GetDag().Nodes_storage_enplacement + "/rootNode/")
	if err != nil {
		panic(fmt.Errorf("UpdateRootNodeFolder could not open folder\nError: %s", err))
	}

	to_add := make([]CRDTDag.EncodedStr, 0)

	for _, file := range files {
		fil, err := os.Open(self.GetDag().Nodes_storage_enplacement + "/rootNode/" + file.Name())
		if err != nil {
			panic(fmt.Errorf("could Not Open RootNode to update rootnodefolder\nerror: %s", err))
		}
		stat, err := fil.Stat()
		if err != nil {
			panic(fmt.Errorf("error in CheckRootNodes, Could not get stat the sub file\nError: %s", err))
		}
		bytesread := make([]byte, stat.Size())
		_, err = fil.Read(bytesread)
		if err != nil {
			panic(fmt.Errorf("error in CheckRootNodes, Could not read the sub file\nError: %s", err))
		}
		err = fil.Close()
		if err != nil {
			panic(fmt.Errorf("error in CheckRootNodes, Could not close the sub file\nError: %s", err))
		}
		if !self.GetDag().IsKnown(bytesread) {
			// separate in 2 folder would be more efficient i think (root note remote and root nodes)
			to_add = append(to_add, CRDTDag.EncodedStr{Str: bytesread})
		}
	}

	self.Merge(to_add)
}
