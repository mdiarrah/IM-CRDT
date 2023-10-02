

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
## Bootstrap
print(bootstrap)
merged_partial=read.csv(bootstrap)
names(merged_partial)[2] ="BootstrapTime"
for (i in 1:n ) {
   print(paste(li[i]))
   dataNode=read.csv(paste(li[i]))
   names(dataNode)[names(dataNode) == "time"] <- paste(paste("Node", i), "Time")
   print(dataNode)
   print(nrow(merged_partial))
   # names(dataNode)[3] = paste(paste("Node", i), "CalculTime")
   merged_partial <- merge(x = merged_partial, y = dataNode,
      by.x = c("CID"),
      by.y = c("CID")
   )
}
print("=======================================merged_partial)==================")
print(merged_partial)
print(colnames(merged_partial))
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

print("==========tmax================")
print(tmax)
names(tmax)[1] = "maxlatency"
merged <- cbind(merged_partial, tmax)


# li <-  append(li, 'V1Res/go_trans_grisou-1/time.csv')
# li <-  append(li, 'V1Res/go_trans_gros-99/time.csv')

# #CRDT ALONE
# CRDTbootstrap = 'V2Res/go_trans_grisou-13/time.csv'
# CRDTli = list()
# CRDTli <-  append(CRDTli, 'V2Res/go_trans_grisou-38/time.csv')
# CRDTli <-  append(CRDTli, 'V2Res/go_trans_grisou-49/time.csv')
# CRDTli <-  append(CRDTli, 'V2Res/go_trans_grisou-50/time.csv')
# CRDTli <-  append(CRDTli, 'V2Res/go_trans_grisou-51/time.csv')
# n=2

# ## CRDTbootstrap

# CRDTmerged_partial=read.csv(CRDTbootstrap)
# names(CRDTmerged_partial)[2] ="BootstrapTime"
# for (i in 1:n ) {
#    dataNode=read.csv(paste(CRDTli[i]))
#    names(dataNode)[2] = paste(paste("Node", i), "Time")
#    CRDTmerged_partial <- merge(x = CRDTmerged_partial, y = dataNode,
#       by.x = c("CID"),
#       by.y = c("CID")
#    )
# }
# tmax = CRDTmerged_partial[paste(paste("Node", 1), "Time")] - CRDTmerged_partial["BootstrapTime"]
# for (i in 1:n ) {
#    latency = CRDTmerged_partial[paste(paste("Node", i), "Time")] -  CRDTmerged_partial["BootstrapTime"]

#    tmax <- pmax(tmax, latency)
#    if (i < n-1) {
#       for (j in i+1:n ) {
#          latency = abs(CRDTmerged_partial[paste(paste("Node", i), "Time")] -  CRDTmerged_partial[paste(paste("Node", j), "Time")])

#          tmax <- pmax(tmax, latency)
#       }
#    }
# }
# names(tmax)[1] =  "maxlatency"
# CRDTmerged <- cbind(CRDTmerged_partial, tmax)
print(merged)


CRDT_IPFS=merged["maxlatency"]
write.csv(mean(CRDT_IPFS$"maxlatency"),file=paste(fileName,'/mean.csv',sep = ""), row.names=TRUE)
# IPFS_ONLY=CRDTmerged["maxlatency"]

sd1=sd(CRDT_IPFS$maxlatency)
# sd2=sd(IPFS_ONLY$maxlatency) 

# pdf("rplot.pdf")
# boxplot(c(CRDT_IPFS,IPFS_ONLY),ylab="time (ms)",names=c("CRDT&IPFS" ,"IPFS_Alone"),col=c("cyan","pink"), 
#    main=paste(paste("MaxLatency of update send, sd:",paste(sd1, " - ")), sd2),horizontal=F)
# dev.off() 
# dev.new(width=5, height=5, unit="in")


################################## RECUP ##########################################
# boxplot(CRDT_IPFS,ylab="time (ms)",names=" Concurrency",col="cyan", 
#    main="CRDT in IPFS MaxLatency on update send",
#    horizontal=F)
# dev.off() 

write.csv(merged,file='Concurrency_1.csv', row.names=TRUE)



# write.csv(CRDTmerged,file='IPFS_ONLY_1.csv', row.names=TRUE)
