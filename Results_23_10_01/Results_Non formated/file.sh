for f in ./*
do
    for x in $f/*Peers
    do
        for xx in $x/*Updater
        do
        echo $xx
            for xxx in $xx/*Updat*
            do
                for xxxx in $xxx/Version*
                do
                    for xxxxx in $xxxx/go_trans_*
                    do
                        find $xxxxx/CRDT_IPFS/node1/ -type f -name "*node*" | parallel rm 
                    done
                done
            done
        done
    done
done
