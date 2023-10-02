package CRDTDag

import (
	CRDT "IPFS_CRDT/Crdt"
	IPFSLink "IPFS_CRDT/ipfsLink"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/path"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

///==============================================================
/// CRDTDag definitions
///==============================================================

type CRDTDag interface {
	Lookup_ToSpecifyType() *CRDT.CRDT
	SendRemoteUpdates()
	Merge(cid EncodedStr)
	GetSys() *IPFSLink.IpfsLink
	GetCRDTManager() *CRDTManager
}
type CRDTManager struct {
	Root_nodes                []EncodedStr
	nodesId                   []([]byte)
	nodesInterface            []*CRDTDagNodeInterface
	Nodes_storage_enplacement string
	pubsubTopic               string
	checkfile                 string
	SubscribedFile            string
	sign                      string
	nodesToAdd_Key            []EncodedStr
	nodesToAdd_value          []*CRDTDagNodeInterface
	nbLineAlreadyWritten      int
	nextNodeName              int
	Sys                       *IPFSLink.IpfsLink
	retrieveMode              bool
}

func (self *CRDTManager) GetSys() *IPFSLink.IpfsLink {
	return self.Sys
}
func Create_CRDTManager(s *IPFSLink.IpfsLink, storage string, bootStrapPeer string) CRDTManager {
	crdt := CRDTManager{
		Root_nodes:                make([]EncodedStr, 0),
		nodesId:                   make([]([]byte), 0),
		nodesInterface:            make([]*CRDTDagNodeInterface, 0),
		Nodes_storage_enplacement: storage,
		SubscribedFile:            "",
		checkfile:                 "",
		sign:                      "",
		pubsubTopic:               "",
		nodesToAdd_Key:            make([]EncodedStr, 0),
		nodesToAdd_value:          make([]*CRDTDagNodeInterface, 0),
		nbLineAlreadyWritten:      0,
		nextNodeName:              0,
		Sys:                       s,
		retrieveMode:              true,
	}
	fmt.Println("storage : ", crdt.Nodes_storage_enplacement)
	crdt.SubscribedFile = crdt.NextFileName()
	if _, err := os.Stat(crdt.SubscribedFile); !errors.Is(err, os.ErrNotExist) {
		os.Remove(crdt.SubscribedFile)
	}
	// if bootStrapPeer == "" {
	// 	crdt.ManageRootNodesConnexion(bootStrapPeer, storage+"rootNode")
	// } else {
	// 	crdt.ManageRootNodesConnexion(bootStrapPeer, storage+"remote")
	// }
	return crdt
}

type EncodedStr struct {
	Str []byte
}

func (self *CRDTManager) GetAllNodes() [][]byte {

	return self.nodesId
}
func (self *CRDTManager) GetAllNodesInterface() []*CRDTDagNodeInterface {
	return self.nodesInterface
}
func (self CRDTManager) EncodeCid(s path.Resolved) EncodedStr {
	b, err := json.Marshal(s.Cid())
	if err != nil {
		panic(fmt.Errorf("Couldn't marshall the path, byte :\nerror : %s", err))
	}
	x := EncodedStr{Str: b}
	return x
}
func (self *CRDTManager) GetNodeFromEncodedCid(s EncodedStr) (files.Node, error) {
	cid := cid.Cid{}
	err := json.Unmarshal(s.Str, &cid)
	if err != nil {
		panic(fmt.Errorf("Couldn't unMarshall the path, byte :%s \nerror : %s", s.Str, err))
	}
	fil, err := IPFSLink.GetIPFS(self.Sys, cid.Bytes())
	if err != nil {
		panic(fmt.Errorf("issue retrieving the IPFS Node :%s", err))
	}
	return fil, err
}

// / @brief Creation of a new empty CRDT Counter in the Operation-based principle
// / @param nodes_stor A free folder where it's possible to write files. The  System will write the nodes here.
// / @param s IPFS System linkin you to the IPFS network, it have to be initialized
func (self *CRDTManager) InitCRDTManager(folderNodeStorage string, s *IPFSLink.IpfsLink, signature int) {

	self.Nodes_storage_enplacement = folderNodeStorage
	self.checkfile = self.NextFileName()
	self.Sys = s
	self.retrieveMode = true
	// self.sign = self.Sys.Hst. Sys HgetID() + "_" + std::to_string(signature)

	self.pubsubTopic = self.sign + "CRDT_" + time.Now().String()
	self.SubscribedFile = folderNodeStorage + "/" + self.pubsubTopic + self.sign + ".data"
	//self.Sys should be subscribed by default
	self.nbLineAlreadyWritten = 0

}
func (self *CRDTManager) NextFileName() string {
	remove_to_save_space := true
	if remove_to_save_space {

		files, err := ioutil.ReadDir(self.Nodes_storage_enplacement)
		if err != nil {
			panic(fmt.Errorf("UpdateRootNodeFolder could not open folder\nError: %s", err))
		}

		for _, file := range files {
			if !file.IsDir() {
				os.Remove(file.Name())
			}
		}
	}
	res := self.Nodes_storage_enplacement + "/node" + strconv.Itoa(self.nextNodeName)
	self.nextNodeName += 1
	return res
}
func (self *CRDTManager) NextRemoteFileName() string {

	res := self.Nodes_storage_enplacement + "/remote/" + strconv.Itoa(self.nextNodeName)
	self.nextNodeName += 1
	return res
}

func (self *CRDTManager) IsKnown(bytes []byte) bool {
	for x := range self.nodesId {
		if string(self.nodesId[x]) == string(bytes) {
			return true
		}
	}
	return false
}

func (self *CRDTManager) UpdateRootNodeFolder() {

	files, err := ioutil.ReadDir(self.Nodes_storage_enplacement + "/rootNode/")
	if err != nil {
		panic(fmt.Errorf("UpdateRootNodeFolder could not open folder\nError: %s", err))
	}

	for _, file := range files {
		if file.Size() > 0 {
			fil, err := os.Open(self.Nodes_storage_enplacement + "/rootNode/" + file.Name())
			if err != nil {
				panic(fmt.Errorf("UPDATE - 1 could Not Open RootNode to update rootnodefolder\nerror: %s", err))
			}
			stat, err := fil.Stat()
			if err != nil {
				panic(fmt.Errorf("UPDATE - error in UpdateRootNode, Could not get stat the sub file\nError: %s", err))
			}
			bytesread := make([]byte, stat.Size())
			_, err = fil.Read(bytesread)
			if err != nil {
				panic(fmt.Errorf("UPDATE - error in UpdateRootNode, Could not read the sub file\nError: %s", err))
			}
			err = fil.Close()
			if err != nil {
				panic(fmt.Errorf("UPDATE - error in UpdateRootNode, Could not close the sub file\nError: %s", err))
			}
			if self.IsKnown(bytesread) {
				// separate in 2 folder would be more efficient i think (root note remote and root nodes)
				err = os.Remove(self.Nodes_storage_enplacement + "/rootNode/" + file.Name())
				if err != nil {
					panic(fmt.Errorf("UPDATE -error in UpdateRootNodeFolder, Could not remove the known file\nError: %s", err))
				}
			}
		}
	}

	for n := range self.Root_nodes {
		fileName := self.Nodes_storage_enplacement + "/rootNode/" + fmt.Sprintf("root%d", self.nextNodeName)
		self.nextNodeName += 1
		fil, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(fmt.Errorf("UPDATE - 2 could Not Open RootNode to update rootnodefolder\nerror: %s", err))
		}
		_, err = fil.Write(self.Root_nodes[n].Str)
		if err != nil {
			panic(fmt.Errorf("could Not write in RootNode to update rootnodefolder\nerror: %s", err))
		}
		err = fil.Close()
		if err != nil {
			panic(fmt.Errorf("could Not Close RootNode to update rootnodefolder\nerror: %s", err))
		}
	}

}

func sendRootNodes(ps *pubsub.PubSub, topic *pubsub.Topic, RootNodeFolder string) {
	files, err := ioutil.ReadDir(RootNodeFolder)
	if err != nil {
		panic(fmt.Errorf("sendRootNodes - Checkupdate could not open folder\nerror: %s", err))
	}

	for _, file := range files {

		fil, err := os.OpenFile(RootNodeFolder+"/remote/"+file.Name(), os.O_RDONLY, os.ModeAppend)
		if err != nil {
			panic(fmt.Errorf("error in sendRootNodes, Could not open the sub file\nError: %s", err))
		}
		stat, err := fil.Stat()
		if err != nil {
			panic(fmt.Errorf("error in sendRootNodes, Could not get stat the sub file\nError: %s", err))
		}
		bytesread := make([]byte, stat.Size())
		_, err = fil.Read(bytesread)
		if err != nil {
			panic(fmt.Errorf("error in sendRootNodes, Could not read the sub file\nError: %s", err))
		}
		err = fil.Close()
		if err != nil {
			panic(fmt.Errorf("error in sendRootNodes, Could not close the sub file\nError: %s", err))
		}
		err = os.Remove(RootNodeFolder + file.Name())
		if err != nil {
			panic(fmt.Errorf("error in sendRootNodes, Could not remove the sub file\nError: %s", err))
		}
		ps.Publish(topic.String(), bytesread)
	}
	ps.Publish(topic.String(), []byte("EOF"))
}

func (self *CRDTManager) ManageRootNodesConnexion(bootstrapPeer string, NodeFolder string) {
	if bootstrapPeer == "" {
		ps2, topic2, sub2 := IPFSLink.SetupPubSub(self.GetSys().Cr.Host, self.GetSys().Cr.Ctx, self.GetSys().Cr.Topic.String()+"_NewConnexions")
		go func() {
			for {
				msg, err := sub2.Next(self.GetSys().Cr.Ctx)
				if err != nil {
					panic(fmt.Errorf("could not read next message on bootstrap peer\nerror: %s", err))
				}
				if msg.GetFrom() != self.GetSys().Cr.Host.ID() {
					sendRootNodes(ps2, topic2, NodeFolder)
				}
			}
		}()
	} else {
		ps2, topic2, sub2 := IPFSLink.SetupPubSub(self.GetSys().Cr.Host, self.GetSys().Cr.Ctx, self.GetSys().Cr.Topic.String()+"_NewConnexions")
		err := ps2.Publish(topic2.String(), []byte("I need the root nodes"))
		if err != nil {
			panic(fmt.Errorf("Could not publish in newConnexion\nError: %s", err))
		}
		go func() {
			i := 0
			msg, err := sub2.Next(self.GetSys().Cr.Ctx)
			if err != nil {
				panic(fmt.Errorf("error in newconnexion pubsub.senxt\nError: %s", err))
			}
			for string(msg.Data) != "EOF" {
				if msg.GetFrom() != self.GetSys().Cr.Host.ID() {
					fil, err := os.OpenFile(NodeFolder+fmt.Sprintf("/rootNode%d", i), os.O_CREATE|os.O_WRONLY, 0755)

					if err != nil {
						panic(fmt.Errorf("could not read next message on bootstrap peer\nerror: %s", err))
					}
					_, err = fil.Write(msg.GetData())
					if err != nil {
						panic(fmt.Errorf("could not write new message from bootstrap peer\nerror: %s", err))
					}
					err = fil.Close()
					if err != nil {
						panic(fmt.Errorf("could not close writen file from bootstrap peer\nerror: %s", err))
					}
				}
			}
			fmt.Println("CANCELLING CONNEXION, I GOT ALL INTERESTING ROOT NODES")
			sub2.Cancel()

		}()
	}
}
func (self *CRDTManager) AddRoot_node(nodeId EncodedStr, node *CRDTDagNodeInterface) {

	self.Root_nodes = append(self.Root_nodes, nodeId)
	self.UpdateRootNodeFolder()

}

func (self *CRDTManager) RemoveRoot_node(nodeId EncodedStr) {

	i := -1
	for x := range self.Root_nodes {
		if string(self.Root_nodes[x].Str) == string(nodeId.Str) {
			i = x
			break
		}
	}
	if i != -1 {
		self.Root_nodes[i] = self.Root_nodes[len(self.Root_nodes)-1]
		self.Root_nodes = self.Root_nodes[:len(self.Root_nodes)-1]
	}
	self.UpdateRootNodeFolder()
}

// /// @brief add the node with the IPFS node ID @node corresponding to @d, assuming no unknown dependance but manage the local root nodes
// /// @param node the node ID to add
// /// @param d The Node itself, with event (data) and direct dependencies (sons)
func (self *CRDTManager) AddNode(node EncodedStr, d *CRDTDagNodeInterface) {

	self.nodesId = append(self.nodesId, node.Str)
	self.nodesInterface = append(self.nodesInterface, d)
	if len((*d).GetDirect_dependency()) > 0 {
		for j := range (*d).GetDirect_dependency() {
			self.RemoveRoot_node((*d).GetDirect_dependency()[j])
		}
	}

	self.AddRoot_node(node, d)
}

// // TODO Check this !!!
func (self *CRDTManager) RemoteAddNodeSuper(cID EncodedStr, newnode *CRDTDagNodeInterface) {

	fmt.Println("=======\nRemote add note\n=======\nPeer:", self.GetSys().Cr.Name, "\nEvent:", (*(*newnode).GetEvent()).ToString(), "\nDirect Dependency:", (*newnode).GetDirect_dependency())
	if self.retrieveMode {
		newNodeFile := ""

		for node := range (*newnode).GetDirect_dependency() {
			if !self.IsKnown((*newnode).GetDirect_dependency()[node].Str) {
				newNodeFile = self.NextFileName()
				if _, err := os.Stat(newNodeFile); !errors.Is(err, os.ErrNotExist) {
					os.Remove(newNodeFile)
				}
				fil, err := self.GetNodeFromEncodedCid((*newnode).GetDirect_dependency()[node])
				if err != nil {
					panic(fmt.Errorf("Remote add failed because %s is malformed\nError: %s", (*newnode).GetDirect_dependency()[node], err))
				}
				files.WriteTo(fil, newNodeFile)
				//Copy of newnode ( to have the type)
				(*newnode).ToFile("/tmp/1")
				var nn *CRDTDagNodeInterface = (*newnode).CreateEmptyNode()
				if err != nil {
					panic(fmt.Errorf("RemoteAddNodeSuper - DeepCopy failed on newnode\nError: %s", err))
				}

				(*nn).FromFile(newNodeFile)

				self.RemoteAddNodeSuper((*newnode).GetDirect_dependency()[node], nn)
			}
		}
		self.AddNode(cID, newnode)
	} else {
		knowAllDependency := true

		// first check if all dependency are resolved
		for node := range (*newnode).GetDirect_dependency() {
			exists := false
			for x := range self.nodesId {
				if string(self.nodesId[x]) == string((*newnode).GetDirect_dependency()[node].Str) {
					exists = true
				}
			}

			// If the key exists
			if exists {
				knowAllDependency = false
				break
			}
		}

		// then Add it weather to :
		// actual node if you have all dependency
		if knowAllDependency {
			stack_EncStr := make([]EncodedStr, 0)
			stack_CRDTDagNode := make([]*CRDTDagNodeInterface, 0)

			stack_EncStr = append(stack_EncStr, cID)
			stack_CRDTDagNode = append(stack_CRDTDagNode, newnode)
			for len(stack_EncStr) > 0 {
				key := stack_EncStr[0]
				actualNode := stack_CRDTDagNode[0]

				//Deleting the first elemnt of the stack ( ^pop)
				stack_CRDTDagNode = stack_CRDTDagNode[1:]
				stack_EncStr = stack_EncStr[:1]
				self.AddNode(key, actualNode)

				for i := range self.nodesToAdd_Key {
					toadd_key := self.nodesToAdd_Key[i]
					toadd_value := self.nodesToAdd_value[i]
					knowAllDependency = true
					for k := range (*toadd_value).GetDirect_dependency() {
						exists := false
						for x := range self.nodesId {
							if string(self.nodesId[x]) == string((*toadd_value).GetDirect_dependency()[k].Str) {
								exists = true
							}
						}
						// If the key exists
						if !exists {
							knowAllDependency = false
							break
						}
					}
					if knowAllDependency {

						stack_CRDTDagNode = append(stack_CRDTDagNode, toadd_value)
						stack_EncStr = append(stack_EncStr, toadd_key)

						self.nodesToAdd_Key[i] = self.nodesToAdd_Key[len(self.nodesToAdd_Key)-1]
						self.nodesToAdd_Key = self.nodesToAdd_Key[:len(self.nodesToAdd_Key)-1]
						self.nodesToAdd_value[i] = self.nodesToAdd_value[len(self.nodesToAdd_value)-1]
						self.nodesToAdd_value = self.nodesToAdd_value[:len(self.nodesToAdd_value)-1]
					}
				}
			}
		} else { // waiting other node to add other people
			self.nodesToAdd_Key = append(self.nodesToAdd_Key, cID)
			self.nodesToAdd_value = append(self.nodesToAdd_value, newnode)
		}

	}

}

func (self *CRDTManager) SendRemoteUpdates() {

	for x := range self.Root_nodes {
		IPFSLink.PubIPFS(self.Sys, self.Root_nodes[x].Str)
	}
}

func CheckForRemoteUpdates(self *CRDTDag, sub *pubsub.Subscription, c context.Context) {
	go func() {
		for {
			msg, err := sub.Next(c)
			if err != nil {

				panic(fmt.Errorf("Check For remote update failed, message not received\nError: %s", err))
			} else if msg.ReceivedFrom == (*self).GetSys().Cr.Host.ID() {
				fmt.Println("Received message from myself")
				continue
			} else {
				fmt.Println("Received message from", msg.ReceivedFrom,
					"data:", string(msg.Data))
				fileName := (*self).GetCRDTManager().NextRemoteFileName()
				if _, err := os.Stat(fileName); !errors.Is(err, os.ErrNotExist) {
					os.Remove(fileName)
				}
				fil, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
				if err != nil {
					panic(fmt.Errorf("error in checkForRemoteUpdate, Could not open the sub file\nError: %s", err))
				}
				_, err = fil.Write(msg.GetData())
				if err != nil {
					panic(fmt.Errorf("error in checkForRemoteUpdate, Could not write the sub file\nError: %s", err))
				}
				err = fil.Close()
				if err != nil {
					panic(fmt.Errorf("error in checkForRemoteUpdate, Could not close the sub file\nError: %s", err))
				}
			}
		}
	}()
}

func (self *CRDTManager) getTopic() string {

	return self.pubsubTopic
}

func (self *CRDTManager) ToString() string {

	str := ""
	str += "nodes : {\n    "
	for s := range self.nodesId {
		str += string(self.nodesId[s]) + " - dd:{"
		for sons := range (*self.nodesInterface[s]).GetDirect_dependency() {
			str += string((*self.nodesInterface[s]).GetDirect_dependency()[sons].Str) + " "
		}
		str += "}\n    "
	}
	str += "}\nRoot_Nodes : {"
	for s := range self.Root_nodes {
		str += string(self.Root_nodes[s].Str) + " "
	}
	str += "}\n"
	return str
}

//TODO : Specify that CID must not be Clear but encoded, so it can be well decrypted by others. ( the only good construction method of Node i Found)
