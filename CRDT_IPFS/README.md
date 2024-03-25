# README #

To run a simple private IPFS network, please follow the steps below:

### Build the binary ###
```
go mod tidy
go build -ldflags "-linkmode 'external' -extldflags '-static'" -o ipfs-crdt
```

### Launch the bootstrap node in a separate directory ###
```
mkdir -p bootstrap && cd bootstrap
cp ../ipfs-crdt .
./ipfs-crdt --bootstrap
```
The bootstrap node will create a json file named peer_info.json in the "demo" directory.
This file contains the bootstrap peer multiaddresses, you need to provide it in the next step (so that new peers can join your network).

### Launch a peer in a separate directory ###
```
mkdir -p peer1 && cd peer1
cp ../ipfs-crdt .
./ipfs-crdt -path ../bootstrap/demo/peer_info.json
```
