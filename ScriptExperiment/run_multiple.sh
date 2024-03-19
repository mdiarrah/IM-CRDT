#!/bin/bash

if [ -z "$1" ]
then
echo "PANIC !!!!!! I NEED an input file representing the used NODES ( similar to OAR_NODE_FILE | sort | uniq )"
else

fileNODE=$1



ARRAY_Repetition=( 1 ) # 4 5 )
ARRAY_NbPeers=( 2 3 4 5 ) # 30 50
ARRAY_UpdatesNb=( 3 5 10 ) #  10 100 
ARRAY_NbPeers_Updating=( 1 2 3 4 5 ) # 30 50


rm advancement

for numeroUNIQUE in "${ARRAY_Repetition[@]}"
do

for nbpeers in "${ARRAY_NbPeers[@]}"
do
for nbpeersUpdating in "${ARRAY_NbPeers_Updating[@]}"
do
for nbupdates in "${ARRAY_UpdatesNb[@]}"
do

if [ $nbpeers -lt $nbpeersUpdating ]
then
echo "$nbpeers < $nbpeersUpdating"
else 

rm "/home/quacher/.ssh/known_hosts"
folder="Results/${nbpeers}Peers/${nbpeersUpdating}Updater/${nbupdates}Updates/Version$numeroUNIQUE"
echo "numeroUNIQUE: $numeroUNIQUE - nbpeers: $nbpeers - nbpeersUpdating: $nbpeersUpdating - nbupdates: $nbupdates" >> advancement
mkdir -p $folder

./run_multipleBIS.sh $nbpeers $nbupdates $nbpeersUpdating  $fileNODE


others=$(cat other)

echo "RETRIEVEDATA"
./RetrieveData.sh resultRetrieve $folder > $folder/Retrieve.log 2>&1
echo "RETRIEVEDATA - THE END"


save=( )

for f in $folder/go_trans_* 
do
     save+=" $f/CRDT_IPFS/node1/time.csv"
done




# Rscript analyseCSV.R $(($nbpeers-1)) "$folder"  ${save[@]}


fi

done
done
done
done


fi
