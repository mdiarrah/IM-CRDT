package Set

import (
	CRDTDag "IPFS_CRDT/CRDTDag"
	CRDT "IPFS_CRDT/Crdt"
	Payload "IPFS_CRDT/Payload"
	IpfsLink "IPFS_CRDT/ipfsLink"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
)

// =======================================================================================
// Payload - OpBased
// =======================================================================================

type Element string
type OpNature int

const (
	ADD OpNature = iota
	REMOVE
)

type Operation struct {
	Elem Element
	Op   OpNature
}

func (self Operation) ToString() string {
	b, err := json.Marshal(self)
	if err != nil {
		panic(fmt.Errorf("Set Operation To string fail to Marshal\nError: %s", err))
	}
	return string(b[:])
}
func (op *Operation) op_from_string(s string) {
	err := json.Unmarshal([]byte(s), op)
	if err != nil {
		panic(fmt.Errorf("Set Operation To string fail to Marshal\nError: %s", err))
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
	b, err := json.Marshal(self)
	if err != nil {
		panic(fmt.Errorf("Set Operation To string fail to Marshal\nError: %s", err))
	}
	return string(b[:])
}
func (self *PayloadOpBased) FromString(s string) {
	err := json.Unmarshal([]byte(s), self)
	if err != nil {
		panic(fmt.Errorf("Set Operation To string fail to Marshal\nError: %s", err))
	}
}

// =======================================================================================
// CRDTSet OpBased
// =======================================================================================

type CRDTSetOpBased struct {
	sys     *IpfsLink.IpfsLink
	added   []string
	removed []string
}

func Create_CRDTSetOpBased(s *IpfsLink.IpfsLink) CRDTSetOpBased {
	return CRDTSetOpBased{
		sys:     s,
		added:   make([]string, 0),
		removed: make([]string, 0),
	}
}
func search(list []string, x string) int {
	for i := 0; i < len(list); i++ {
		if list[i] == x {
			return i
		}
	}
	return -1
}
func (self *CRDTSetOpBased) Add(x string) {
	if search(self.added, x) == -1 {
		self.added = append(self.added, x)
	}
}

func (self *CRDTSetOpBased) Remove(x string) {
	if search(self.removed, x) == -1 {
		self.removed = append(self.removed, x)
	}
}

func (self *CRDTSetOpBased) Lookup() []string {
	i := make([]string, 0)
	fmt.Println("size", len(self.added))
	for x := range self.added {
		if search(self.removed, self.added[x]) == -1 {
			i = append(i, self.added[x])
			i = append(i, ",")
		}
	}

	return i
}

func (self *CRDTSetOpBased) ToFile(file string) {

	b, err := json.Marshal(self)
	if err != nil {
		panic(fmt.Errorf("CRDTDagNode - ToFile Could not Marshall %s\nError: %s", file, err))
	}
	f, err := os.Create(file)
	if err != nil {
		panic(fmt.Errorf("CRDTDagNode - ToFile Could not Create the file %s\nError: %s", file, err))
	}
	f.Write(b)
	err = f.Close()
	if err != nil {
		panic(fmt.Errorf("CRDTDagNode - ToFile Could not Write to the file %s\nError: %s", file, err))
	}
}

// =======================================================================================
// CRDTSetDagNode OpBased
// =======================================================================================

type CRDTSetOpBasedDagNode struct {
	DagNode CRDTDag.CRDTDagNode
}

func (self *CRDTSetOpBasedDagNode) FromFile(fil string) {
	var pl Payload.Payload = &PayloadOpBased{}
	self.DagNode.CreateNodeFromFile(fil, &pl)
}

func (self *CRDTSetOpBasedDagNode) GetDirect_dependency() []CRDTDag.EncodedStr {

	return self.DagNode.DirectDependency
}

func (self *CRDTSetOpBasedDagNode) ToFile(file string) {

	self.DagNode.ToFile(file)
}
func (self *CRDTSetOpBasedDagNode) GetEvent() *Payload.Payload {

	return self.DagNode.Event
}
func (self *CRDTSetOpBasedDagNode) GetPiD() string {

	return self.DagNode.PID
}
func (self *CRDTSetOpBasedDagNode) CreateEmptyNode() *CRDTDag.CRDTDagNodeInterface {
	n := CreateDagNode(Operation{}, "")
	var node CRDTDag.CRDTDagNodeInterface = &n
	return &node
}
func CreateDagNode(o Operation, id string) CRDTSetOpBasedDagNode {
	var pl Payload.Payload = &PayloadOpBased{Op: o, Id: id}
	slic := make([]CRDTDag.EncodedStr, 0)
	return CRDTSetOpBasedDagNode{
		DagNode: CRDTDag.CRDTDagNode{
			Event:            &pl,
			PID:              id,
			DirectDependency: slic,
		},
	}
}

// =======================================================================================
// CRDTSetDag OpBased
// =======================================================================================

type CRDTSetOpBasedDag struct {
	dag CRDTDag.CRDTManager
}

func (self *CRDTSetOpBasedDag) GetDag() *CRDTDag.CRDTManager {

	return &self.dag
}
func (self *CRDTSetOpBasedDag) SendRemoteUpdates() {

	self.dag.SendRemoteUpdates()
}
func (self *CRDTSetOpBasedDag) GetCRDTManager() *CRDTDag.CRDTManager {

	return &self.dag
}
func (self *CRDTSetOpBasedDag) IsKnown(cid CRDTDag.EncodedStr) bool {

	find := false
	for x := range self.dag.GetAllNodes() {
		if string(self.dag.GetAllNodes()[x]) == string(cid.Str) {
			find = true
			break
		}
	}
	return find
}
func (self *CRDTSetOpBasedDag) Merge(cid CRDTDag.EncodedStr) {

	find := self.IsKnown(cid)
	if !find {
		fil, err := self.dag.GetNodeFromEncodedCid(cid)
		if err != nil {
			panic(fmt.Errorf("could not retrieve the node %s , error :%s", cid.Str, err))
		}
		fstr := self.dag.NextFileName()
		if _, err := os.Stat(fstr); !errors.Is(err, os.ErrNotExist) {
			os.Remove(fstr)
		}
		files.WriteTo(fil, fstr)
		n := CreateDagNode(Operation{}, "")
		n.FromFile(fstr)

		self.remoteAddNode(cid, n)
	}
}
func (self *CRDTSetOpBasedDag) remoteAddNode(cID CRDTDag.EncodedStr, newnode CRDTSetOpBasedDagNode) {
	var pl CRDTDag.CRDTDagNodeInterface = &newnode
	self.dag.RemoteAddNodeSuper(cID, &pl)
}

func (self *CRDTSetOpBasedDag) Add(x string) string {
	newNode := CreateDagNode(Operation{Elem: Element(x), Op: ADD}, self.GetSys().Hst.ID().Pretty())
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
		panic(fmt.Errorf("ERROR INCREMENT CRDTSetOpBasedDag, could not read file\nerror: %s", err))
	}
	path, err := IpfsLink.AddIPFS(self.dag.Sys, bytes)
	if err != nil {
		panic(fmt.Errorf("CRDTSetOpBasedDag Increment, could not add the file to IFPS\nerror: %s", err))
	}

	encodedCid := self.dag.EncodeCid(path)
	c := cid.Cid{}
	err = json.Unmarshal(encodedCid.Str, &c)
	if err != nil {
		panic(fmt.Errorf("CRDTSetOpBasedDag Increment, could not UnMarshal\nerror: %s", err))
	}

	// fmt.Println("encodedCid Increment :", c.String())
	var pl CRDTDag.CRDTDagNodeInterface = &newNode

	self.dag.AddNode(encodedCid, &pl) // TODOSetCrdt Complete Node interface

	self.SendRemoteUpdates()
	return c.String()
}
func (self *CRDTSetOpBasedDag) Remove(x string) string {

	newNode := CreateDagNode(Operation{Elem: Element(x), Op: REMOVE}, self.GetSys().Hst.ID().Pretty())
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
		panic(fmt.Errorf("ERROR INCREMENT CRDTSetOpBasedDag, could not read file\nerror: %s", err))
	}
	path, err := IpfsLink.AddIPFS(self.dag.Sys, bytes)
	if err != nil {
		panic(fmt.Errorf("CRDTSetOpBasedDag Decrement, could not add the file to IFPS\nerror: %s", err))
	}

	encodedCid := self.dag.EncodeCid(path)
	c := cid.Cid{}
	err = json.Unmarshal(encodedCid.Str, &c)
	if err != nil {
		panic(fmt.Errorf("CRDTSetOpBasedDag Increment, could not UnMarshal\nerror: %s", err))
	}

	// _, c, _ := cid.CidFromBytes(encodedCid.Str)
	// fmt.Println("encodedCid Decrement :", c.String())
	var pl CRDTDag.CRDTDagNodeInterface = &newNode
	self.dag.AddNode(encodedCid, &pl)
	self.SendRemoteUpdates()

	return c.String()
}

func Create_CRDTSetOpBasedDag(sys *IpfsLink.IpfsLink, storage_emplacement string, bootStrapPeer string) CRDTSetOpBasedDag {
	man := CRDTDag.Create_CRDTManager(sys, storage_emplacement, bootStrapPeer)
	crdtSet := CRDTSetOpBasedDag{dag: man}
	if bootStrapPeer == "" {
		x, err := os.ReadFile("initial_value")
		if err != nil {
			panic(fmt.Errorf("Could not read initial_value, error : %s", err))
		}
		newNode := CreateDagNode(Operation{Elem: Element(x), Op: ADD}, crdtSet.GetSys().Hst.ID().Pretty())

		strFile := crdtSet.dag.NextFileName()
		if _, err := os.Stat(strFile); !errors.Is(err, os.ErrNotExist) {
			os.Remove(strFile)
		}
		newNode.ToFile(strFile)
		bytes, err := os.ReadFile(strFile)
		if err != nil {
			panic(fmt.Errorf("ERROR INCREMENT CRDTSetOpBasedDag, could not read file\nerror: %s", err))
		}
		path, err := IpfsLink.AddIPFS(crdtSet.dag.Sys, bytes)
		if err != nil {
			panic(fmt.Errorf("CRDTSetOpBasedDag Increment, could not add the file to IFPS\nerror: %s", err))
		}

		encodedCid := crdtSet.dag.EncodeCid(path)
		c := cid.Cid{}
		err = json.Unmarshal(encodedCid.Str, &c)
		if err != nil {
			panic(fmt.Errorf("CRDTSetOpBasedDag Increment, could not UnMarshal\nerror: %s", err))
		}

		// fmt.Println("encodedCid Increment :", c.String())
		var pl1 CRDTDag.CRDTDagNodeInterface = &newNode

		crdtSet.dag.AddNode(encodedCid, &pl1) // TODOSetCrdt Complete Node interface
	}
	var pl CRDTDag.CRDTDag = &crdtSet

	CRDTDag.CheckForRemoteUpdates(&pl, sys.Cr.Sub, man.Sys.Ctx)
	return crdtSet
}

func (self *CRDTSetOpBasedDag) GetSys() *IpfsLink.IpfsLink {

	return self.dag.Sys
}

func (self *CRDTSetOpBasedDag) Lookup_ToSpecifyType() *CRDT.CRDT {

	crdt := CRDTSetOpBased{
		sys:     self.GetSys(),
		added:   make([]string, 0),
		removed: make([]string, 0),
	}
	for x := range self.dag.GetAllNodes() {
		node := self.dag.GetAllNodesInterface()[x]
		if (*(*node).GetEvent()).(*PayloadOpBased).Op.Op == ADD {
			// fmt.Println("add")
			crdt.Add(string((*(*node).GetEvent()).(*PayloadOpBased).Op.Elem))
		} else {
			// fmt.Println("remove")
			crdt.Remove(string((*(*node).GetEvent()).(*PayloadOpBased).Op.Elem))
		}
	}
	var pl CRDT.CRDT = &crdt
	return &pl
}
func (self *CRDTSetOpBasedDag) Lookup() CRDTSetOpBased {

	// crdt := self.Lookup_ToSpecifyType()
	// var pl CRDTDag.CRDTDag = &crdtSet
	return *(*self.Lookup_ToSpecifyType()).(*CRDTSetOpBased)
}

type Tuple struct {
	Cid                string
	IntegrityCheckTime int
	CalculTime         int
}

func (self *CRDTSetOpBasedDag) CheckUpdate() []Tuple {
	files, err := ioutil.ReadDir(self.GetDag().Nodes_storage_enplacement + "/remote")
	received := make([]Tuple, 0)
	if err != nil {
		fmt.Printf("CheckUpdate - Checkupdate could not open folder\nerror: %s\n", err)
	} else {
		for _, file := range files {
			if file.Size() > 0 {
				ti := time.Now()
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

				// fmt.Println("stat.size :", stat.Size(), "read :", n)
				if int64(n) != stat.Size() {
					panic(fmt.Errorf("error in checkupdate, Could not read entirely the sub file\nError: read %d byte unstead of %d", n, stat.Size()))
				}
				err = fil.Close()
				if err != nil {
					panic(fmt.Errorf("error in checkupdate, Could not close the sub file\nError: %s", err))
				}
				timeToCompute := time.Since(ti).Nanoseconds()
				ti = time.Now()
				if !self.IsKnown(CRDTDag.EncodedStr{Str: bytesread}) {
					self.Merge(CRDTDag.EncodedStr{Str: bytesread})
				}
				checktime := time.Since(ti).Nanoseconds()
				ti = time.Now()
				s := cid.Cid{}
				json.Unmarshal(bytesread, &s)

				err = os.Remove(self.GetDag().Nodes_storage_enplacement + "/remote/" + file.Name())
				if err != nil || errors.Is(err, os.ErrNotExist) {
					panic(fmt.Errorf("error in checkupdate, Could not remove the sub file\nError: %s", err))
				}
				timeToCompute += time.Since(ti).Nanoseconds()
				received = append(received, Tuple{Cid: s.String(), IntegrityCheckTime: int(checktime), CalculTime: int(timeToCompute)})

			} else {
				fmt.Printf("FILE SIZE NULL")
			}
		}
		// fmt.Println("calling UpdateRootNodeFolder")
		self.GetDag().UpdateRootNodeFolder()
	}
	return received
}

// Next function must be useless
func (self *CRDTSetOpBasedDag) CheckRootNodes() {
	files, err := ioutil.ReadDir(self.GetDag().Nodes_storage_enplacement + "/rootNode/")
	if err != nil {
		panic(fmt.Errorf("UpdateRootNodeFolder could not open folder\nError: %s", err))
	}

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
			self.Merge(CRDTDag.EncodedStr{Str: bytesread})
		}
	}

}
