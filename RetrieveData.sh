#!/bin/bash


folder=$2

rm -r "StudyResult/$folder"
mkdir "StudyResult/$folder"

tar -xf $1 -C StudyResult/$folder

for filename in ./StudyResult/$folder/resultRetrieve/*.tar.gz; do

    EXTENSION=`echo "$filename" | cut -d'.' -f2`
    EXTENSION=`echo "$EXTENSION" | cut -d'/' -f5`
    mkdir "./StudyResult/$folder/$EXTENSION"
    tar -xf $filename -C ./StudyResult/$folder/$EXTENSION
    mv ./StudyResult/$folder/$EXTENSION/go_transcription2/node1/* ./StudyResult/$folder/$EXTENSION
done

#tar -xvf go_trans.tar.gz > tar.log 
