

###################################################################################################
###################################################################################################
#usage : 
# Rscript analyseCSV.R ${N} ${IDENTIFICATION_NUMBER} ${BOOTSTRAP} ${NODE1} ${NODE2} ... ${NODEN} 
###################################################################################################
###################################################################################################


#CRDT+IPFS
args <- commandArgs()
number = args[6]
fileName = args[7]
bootstrap = args[8]

li = list('')
li[1]  = args[9]
for (i in 1:(number)) {
li <- append(li, args[9+i])
}
n=strtoi(number)

merged_partial=read.csv(bootstrap)
names(merged_partial)[2] ="BootstrapTime"
for (i in 1:n ) {
   print(paste(li[i]))
   dataNode=read.csv(paste(li[i]))
   names(dataNode)[names(dataNode) == "time"] <- paste(paste("Node", i), "Time")
   # names(dataNode)[3] = paste(paste("Node", i), "CalculTime")
   merged_partial <- merge(x = merged_partial, y = dataNode,
      by.x = c("CID"),
      by.y = c("CID")
   )
}

tmax = abs(merged_partial["Node 1 Time"] - merged_partial["BootstrapTime"])
for (i in 1:n ) {
   latency = abs(merged_partial[paste(paste("Node", i), "Time")] -  merged_partial["BootstrapTime"])

   tmax <- pmax(tmax, latency)

   if (i < n - 1) {
      for (j in i+1:n-1 ) {
         if (j < n) {
         latency = abs(merged_partial[paste(paste("Node", i), "Time")] - merged_partial[paste(paste("Node", j), "Time")])

         tmax <- pmax(tmax, latency)
         }
      }
   }
}


names(tmax)[1] = "maxlatency"
merged <- cbind(merged_partial, tmax)




CRDT_IPFS=merged["maxlatency"]
write.csv(mean(CRDT_IPFS$"maxlatency"),file=paste(fileName,'/mean.csv',sep = ""), row.names=TRUE)

sd1=sd(CRDT_IPFS$maxlatency)

write.csv(merged,file='Concurrency_1.csv', row.names=TRUE)

