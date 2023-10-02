#!/bin/bash


echo "/**"
echo " * Beggining the CRDT IPFS setup"
echo " */"

USER_LOGIN_NAME=$(id -un)
USER_GROUP_ID=$(id -g)
USER_GROUP_NAME=$(id -gn)

DATE=$(date +%s)

TMP_DIR=/tmp/$DATE'-'$$'-CRDTIPFS'

echo $TMP_DIR

let NHOST=($(cat $OAR_FILE_NODES | uniq | wc -l)-1)
echo "numHost : $NHOST"
SLAVES=$(cat $OAR_FILE_NODES | uniq | tail -n $NHOST)
MASTER=$(head -n 1 $OAR_FILE_NODES)
ALL_NODES=$(cat $OAR_FILE_NODES | uniq)

# outside
# tar czvf go_trans.tar.gz go_transcription2
# scp go_trans.tar.gz Nancy.g5k:~/go_trans.tar.gz
# ssh Nancy.g5k

# inside
# oarsub -I -l host=1,walltime=1:45 -t deploy
# kadeploy3 -a IPFS_CRDT_Environnement_GO.yaml
# scp  ./go_trans.tar.gz root@grisou-31:~/go_trans.tar.gz

# ssh root@grisou-31

# tar -xvf go_trans.tar.gz
# cd go_transcription2
# go build

mkdir $TMP_DIR

cat > $TMP_DIR/masters << EOF
$MASTER
EOF

cat > $TMP_DIR/slaves << EOF
$SLAVES
EOF


kadeploy3 -a CRDT_IPFS.yaml

echo "Building the  GO implementation"

for SLAVE in $ALL_NODES
do
scp $1 root@$SLAVE:~/$1
ssh root@$SLAVE "sh -c 'tar -xvf go_trans.tar.gz > tar.log && cd go_transcription2 && /usr/local/go/bin/go build > build.log'"
done

echo "running the bootstrap in ${MASTER}"
ssh root@$MASTER "rm  go_transcription2/ID"
ssh root@$MASTER "sh -c 'cd go_transcription2 && ./IPFS_CRDT --mode BootStrap > out.log & '" &
echo "done, now i sleep"
sleep 10s
echo "done sleeping"
BOOTSTRAPID=$(ssh root@$MASTER "sh -c 'cat ./go_transcription2/ID'")

echo "running the lisnteners in $SLAVES with bootstrapID: $BOOTSTRAPID"

for SLAVE in $SLAVES
do

ssh root@$SLAVE "rm -rf go_transcription2/FIRST"
ssh root@$SLAVE "sh -c 'cd go_transcription2 && ./IPFS_CRDT --mode update --ni ${BOOTSTRAPID} > out.log &'"&
done


rm -rf $TMP_DIR/
