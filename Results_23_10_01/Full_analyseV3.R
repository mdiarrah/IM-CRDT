
###################################################################################################
###################################################################################################
#usage : 
# Rscript analyseCSV.R  ${Output_folder} ${file_of_input_files} 
# input file :
# nb_peers,nb_update,Version,Mean,system
# 5,100,1,360,IPFS+CRDT
# 10,1000,1,450,IPFS_ALONE
# ...,...,...,...
###################################################################################################
###################################################################################################

library(conflicted)  
# library(dplyr)
library(ggplot2)
library(tidyverse)

conflict_prefer("filter", "dplyr")
conflict_prefer("lag", "dplyr")

# IM-CRDT
# 2 Peers

data_Frame_IM_CRDT = data.frame(CID=character(0),maxlatency=numeric(0),mean_latency=numeric(0),mean_time_retrieve=numeric(0),mean_time_compute=numeric(0),time_add_IPFS=numeric(0),mean_time_pubsub=numeric(0),numUpdates=numeric(0),numberPeers=numeric(0),numberPeerUpdating=numeric(0),System=character(0))

for (nb_peers in c(2, 5, 10, 20, 50) )
{
 for (nb_Peers_Updating in c(1, 2, 5, 10, 20) )
  {
    if (nb_peers >= nb_Peers_Updating)
    {
      for (nb_Updates in c(10, 100, 1000) )
      {
        file=paste("DATA_Experience/IM-CRDT/output_",nb_peers,"Peers_",nb_Peers_Updating,"Updater_",nb_Updates,"Updates.csv", sep = "")
        print(file)
        a = read.csv(file)
        b <- a %>%
          add_column(numUpdates = nb_Updates) %>%
          add_column(numberPeers = nb_peers) %>%
          add_column(numberPeerUpdating = nb_Peers_Updating) %>%
          add_column(System = "IM-CRDT")
        data_Frame_IM_CRDT= rbind(data_Frame_IM_CRDT, b)

      }
    }
  } 
}


write.table(data_Frame_IM_CRDT, "DATA_Experience/totalDATAFRAME_IM-CRDT.csv")

data_Frame_IPFS_Alone = data.frame(CID=character(0),maxlatency=numeric(0),mean_latency=numeric(0),mean_timeRetrieve=numeric(0),mean_timeSend=numeric(0),mean_time_pubsub=numeric(0),numUpdates=numeric(0),numberPeers=numeric(0),numberPeerUpdating=numeric(0),System=character(0))

for (nb_peers in c(2, 5, 10, 20, 50) )
{
    for (nb_Updates in c(10, 100, 1000) )
    {    
        file=paste("DATA_Experience/IPFS_Alone/output_",nb_peers,"Peers_1Updater_",nb_Updates,"Updates.csv", sep = "")
        print(file)
        a = read.csv(file)
        b <- a %>%
          add_column(numUpdates = nb_Updates) %>%
          add_column(numberPeers = nb_peers) %>%
          add_column(numberPeerUpdating = 1) %>%
          add_column(System = "IPFS Alone")
        data_Frame_IPFS_Alone= rbind(data_Frame_IPFS_Alone, b)

    }
}


write.table(data_Frame_IPFS_Alone, "DATA_Experience/totalDATAFRAME_IPFS_Alone.csv")




##############################- Comparison between IM-CRDT and IPFS Alone -##############################
# 10 Updates


data_IM_CRDT = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 1)

data_IPFS_Alone = filter(data_Frame_IPFS_Alone, System == "IPFS Alone"  & numberPeerUpdating == 1)


data_IM_CRDT= data.frame(maxlatency=data_IM_CRDT$"maxlatency", System=data_IM_CRDT$"System", numberPeers=data_IM_CRDT$"numberPeers", numUpdates=data_IM_CRDT$"numUpdates")

data_IPFS_Alone= data.frame(maxlatency=data_IPFS_Alone$"maxlatency", System=data_IPFS_Alone$"System", numberPeers=data_IPFS_Alone$"numberPeers", numUpdates=data_IPFS_Alone$"numUpdates")


data=rbind(data_IM_CRDT,data_IPFS_Alone)
data$numberPeers = as.character(data$numberPeers)

# grouped boxplot

p <- data %>%
  mutate(numberPeers = fct_relevel(numberPeers, "2", "5", "10", "20", "50")) %>%
  ggplot( aes(x=numberPeers, y=maxlatency, fill=System)) +
    theme(text = element_text(size = 20)) +
    geom_boxplot(outlier.shape = NA) +
    facet_wrap(~numUpdates) +
    xlab("Number of replicas")+
    ylab("Maximum latency (ms)") +
    scale_y_continuous(limits = c(100,600)) 

#     geom_boxplot(outlier.shape = NA, fill=System) 

plot(p)

# ggplot(data_IPFS_Alone, aes(x=numberPeers, y=maxlatency)) + 
#     geom_boxplot(outlier.shape = NA, fill=System) 




##############################- Test of scalability of IM-CRDT -##############################



# grouped boxplot
data_Frame_IM_CRDT$numberPeers = as.character(data_Frame_IM_CRDT$numberPeers)
data_Frame_IM_CRDT$numberPeerUpdating = as.character(data_Frame_IM_CRDT$numberPeerUpdating)

p <- data_Frame_IM_CRDT %>%
  mutate(numberPeerUpdating = fct_relevel(numberPeerUpdating, "1" , "2", "5", "10", "20")) %>%
  mutate(numberPeers = fct_relevel(numberPeers, "2" , "5", "10", "20", "50")) %>%
  ggplot( aes(x=numberPeerUpdating, y=maxlatency, fill=numberPeers)) +
    theme(text = element_text(size = 20)) +
    geom_boxplot(outlier.shape = NA) +
    facet_wrap(~numUpdates) +
    labs(fill = "replicas") + 
    xlab("Number of peers updating")+
    ylab("Maximum latency (ms)") +
    scale_y_log10()


plot(p)


##############################- evolution of different times -##############################





# create a dataset
specie <- c(rep(1 , 4) , rep(2 , 4) , rep(3 , 4) , rep(4 , 4), rep(5, 4) )
condition <- rep(c("mean_time_retrieve" , "mean_time_compute" , "time_add_IPFS", "mean_time_pubsub") , 5)


CSVFrameCRDT_IPFS    = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 1)
CSVFrameCRDT_IPFS_2  = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 2)
CSVFrameCRDT_IPFS_5  = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 5)
CSVFrameCRDT_IPFS_10 = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 10)
CSVFrameCRDT_IPFS_20 = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 20)

print("CSVFrameCRDT_IPFS_20")
# print(CSVFrameCRDT_IPFS_20)

value = abs(c(mean(CSVFrameCRDT_IPFS$"mean_time_retrieve") /  1000000,mean(CSVFrameCRDT_IPFS$"mean_time_compute")   / 1000000 ,mean(CSVFrameCRDT_IPFS$"time_add_IPFS")    / 1000000,mean(CSVFrameCRDT_IPFS$"mean_time_pubsub")   ,
    mean(CSVFrameCRDT_IPFS_2$"mean_time_retrieve") /  1000000,mean(CSVFrameCRDT_IPFS_2$"mean_time_compute") / 1000000 ,mean(CSVFrameCRDT_IPFS_2$"time_add_IPFS")  / 1000000,mean(CSVFrameCRDT_IPFS_2$"mean_time_pubsub") ,
    mean(CSVFrameCRDT_IPFS_5$"mean_time_retrieve") /  1000000,mean(CSVFrameCRDT_IPFS_5$"mean_time_compute") / 1000000 ,mean(CSVFrameCRDT_IPFS_5$"time_add_IPFS")  / 1000000,mean(CSVFrameCRDT_IPFS_5$"mean_time_pubsub") ,
    mean(CSVFrameCRDT_IPFS_10$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS_10$"mean_time_compute" / 1000000),mean(CSVFrameCRDT_IPFS_10$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS_10$"mean_time_pubsub"),
    mean(CSVFrameCRDT_IPFS_20$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS_20$"mean_time_compute" / 1000000),mean(CSVFrameCRDT_IPFS_20$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS_20$"mean_time_pubsub", na.rm=TRUE)))
data <- data.frame(specie,condition,value)
# print(data)

print(CSVFrameCRDT_IPFS_20$"mean_time_pubsub")




ggplot(data, aes(fill=condition, y=value, x=specie)) + 
    theme(text = element_text(size = 20)) +
    scale_x_continuous(,breaks = seq(1, 5, 1), labels = c(1,2,5,10,20))+
    xlab("Number of peers updating")+
    ylab("time (ms)") +
    scale_fill_discrete(labels = c("Compute", "Pubsub", "Retrieve", "Add IPFS")) +
    labs(fill = "Step") + 
    geom_bar(position="stack", stat="identity") + 
    scale_y_log10()


specieCompute <- c(rep(1 , 1) , rep(2 , 1) , rep(3 , 1) , rep(4 , 1), rep(5, 1) )
conditionCompute <- rep(c("mean_time_compute" ), 5)
valueCompute = abs(c(mean(CSVFrameCRDT_IPFS$"mean_time_compute") ,
    mean(CSVFrameCRDT_IPFS_2$"mean_time_compute"),
    mean(CSVFrameCRDT_IPFS_5$"mean_time_compute"),
    mean(CSVFrameCRDT_IPFS_10$"mean_time_compute"),
    mean(CSVFrameCRDT_IPFS_20$"mean_time_compute")))
dataCompute<- data.frame(specieCompute,conditionCompute,valueCompute)
print(dataCompute)


ggplot(dataCompute, aes(y=valueCompute, x=specieCompute)) + 
    theme(text = element_text(size = 20)) +
    scale_x_continuous(,breaks = seq(1, 5, 1), labels = c(1,2,5,10,20))+
    xlab("Number of peers updating")+
    ylab("Compute time (ns)") +
    geom_line() +
    geom_point()





##############################- evolution of retrieve times -##############################


# colnames(data_IPFS_Alone)[which(names(df) == "mean_timeRetrieve")] <- "mean_time_retrieve"



data_IM_CRDT = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 1)

data_IPFS_Alone = filter(data_Frame_IPFS_Alone, System == "IPFS Alone"  & numberPeerUpdating == 1)


data_IM_CRDT= data.frame(mean_time_retrieve=data_IM_CRDT$"mean_time_retrieve"/1000000, System=data_IM_CRDT$"System", numberPeers=data_IM_CRDT$"numberPeers", numUpdates=data_IM_CRDT$"numUpdates")

data_IPFS_Alone= data.frame(mean_time_retrieve=data_IPFS_Alone$"mean_timeRetrieve"/1000000, System=data_IPFS_Alone$"System", numberPeers=data_IPFS_Alone$"numberPeers", numUpdates=data_IPFS_Alone$"numUpdates")


data=rbind(data_IM_CRDT,data_IPFS_Alone)
data$numberPeers = as.character(data$numberPeers)


print(data)
# grouped boxplot

p <- data %>%
  mutate(numberPeers = fct_relevel(numberPeers, "2", "5", "10", "20", "50")) %>%
  ggplot( aes(x=numberPeers, y=mean_time_retrieve, fill=System)) +
    theme(text = element_text(size = 20)) +
    geom_boxplot(outlier.shape = NA) +
    facet_wrap(~numUpdates) +
    xlab("Number of replicas")+
    ylab("Mean Time to retrieve (ms)") +
    scale_y_continuous(limits = c(30,275)) 

#     geom_boxplot(outlier.shape = NA, fill=System) 

plot(p)





# data_IM_CRDT = filter(data_Frame_IM_CRDT, System == "IM-CRDT"  & numberPeerUpdating == 1)

# data_IPFS_Alone = filter(data_Frame_IPFS_Alone, System == "IPFS Alone"  & numberPeerUpdating == 1)


# data_IM_CRDT= data.frame(maxlatency=data_IM_CRDT$"maxlatency", System=data_IM_CRDT$"System", numberPeers=data_IM_CRDT$"numberPeers", numUpdates=data_IM_CRDT$"numUpdates")

# data_IPFS_Alone= data.frame(maxlatency=data_IPFS_Alone$"maxlatency", System=data_IPFS_Alone$"System", numberPeers=data_IPFS_Alone$"numberPeers", numUpdates=data_IPFS_Alone$"numUpdates")


# p <- data_Frame_IM_CRDT %>%
#   mutate(numberPeerUpdating = fct_relevel(numberPeerUpdating, "1" , "2", "5", "10", "20")) %>%
#   ggplot( aes(x=numberPeerUpdating, y=maxlatency, fill=numberPeers)) +
#     geom_bar(position="dodge", stat="identity") +
#     facet_wrap(~numUpdates) +
#     scale_y_log10()


# plot(p)


##############################- Comparison between IM-CRDT and IPFS Alone -##############################





# create a dataset
# specie <- c(rep(1 , 4) , rep(2 , 4) , rep(3 , 4) , rep(4 , 4), rep(5, 4) )
# condition <- rep(c("mean_time_retrieve" , "mean_time_compute" , "time_add_IPFS", "mean_time_pubsub") , 5)
# value <- 
# abs(c(mean(CSVFrameCRDT_IPFS$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS$"mean_time_compute") / 1000000,mean(CSVFrameCRDT_IPFS$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS$"mean_time_pubsub"),
# mean(CSVFrameCRDT_IPFS_2$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS_2$"mean_time_compute") / 1000000,mean(CSVFrameCRDT_IPFS_2$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS_2$"mean_time_pubsub"),
# mean(CSVFrameCRDT_IPFS_5$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS_5$"mean_time_compute") / 1000000,mean(CSVFrameCRDT_IPFS_5$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS_5$"mean_time_pubsub"),
# mean(CSVFrameCRDT_IPFS_10$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS_10$"mean_time_compute" / 1000000),mean(CSVFrameCRDT_IPFS_10$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS_10$"mean_time_pubsub"),
# mean(CSVFrameCRDT_IPFS_20$"mean_time_retrieve") / 1000000,mean(CSVFrameCRDT_IPFS_20$"mean_time_compute" / 1000000),mean(CSVFrameCRDT_IPFS_20$"time_add_IPFS") / 1000000,mean(CSVFrameCRDT_IPFS_20$"mean_time_pubsub")))
# data <- data.frame(specie,condition,value)
 

# print(paste("mean_time_retrieve", mean(CSVFrameCRDT_IPFS_10$"mean_time_retrieve")))
# print(paste("mean_time_compute", mean(CSVFrameCRDT_IPFS_10$"mean_time_compute")))
# print(paste("time_add_IPFS", mean(CSVFrameCRDT_IPFS_10$"time_add_IPFS")))
# print(paste("mean_time_pubsub", mean(CSVFrameCRDT_IPFS_10$"mean_time_pubsub")))

# pdf("AnalyseCommunicationTime.pdf")
# ggplot(data, aes(fill=condition, y=value, x=specie)) + 
#     scale_x_continuous(,breaks = seq(1, 5, 1), labels = c(1,2,5,10,20))+
#     xlab("Number of peers updating")+
#     ylab("Percentage of time spend on the step") +
#     scale_fill_discrete(labels = c("Compute", "Pubsub", "Retrieve", "Add IPFS")) +
#     ggtitle("Comparison of communication and computation time \n in function of the number of updates") +
#     labs(fill = "Step") + 
#     geom_bar(position="fill", stat="identity")













#USED PART /////////////////////////////////////////////////////////////
# pdf("comparison,10Updates Output.pdf")
# # CSVFrame1=rbind(CSVFrame1,c(1,1))
# print("0")
# boxplot(cbind(CSVFrameCRDT_IPFS$"maxlatency", CSVFrameCRDT_IPFS_5$"maxlatency", CSVFrameCRDT_IPFS_10$"maxlatency", CSVFrameCRDT_IPFS_20$"maxlatency", CSVFrameCRDT_IPFS_50$"maxlatency"), ylim=c(0,700),outline=FALSE,main="IM-CRDT latency analysis 1 Updater, 10Updates",at= c(2,5,10,20,50),outlier.shape = NA,  names=c("2Peers", "5Peers",  "10Peers", "20Peers", "50Peers"), xlab="Number of peers",ylab="Maximum latency per update (ms)", boxfill = NA, border = NA) #invisible boxes - only axes and plot area

# print("1")
# #print(CSVFrameCRDT_IPFS_100$"maxlatency")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS$"maxlatency", xaxt="n",at=c(2),add = TRUE, boxwex = 1, outlier.shape = NA, data=mtcars, names=c("1"),col=colors[1]) # at exprime la position en x

# print("2")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS_5$"maxlatency", xaxt="n",at=c(5),add = TRUE,boxwex = 1, data=mtcars, names=c("2"),col=colors[1]) # at exprime la position en x

# print("3")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS_10$"maxlatency", xaxt="n",at=c(10),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("5"),col=colors[1]) # at exprime la position en x

# print("4")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS_20$"maxlatency", xaxt="n",at=c(20),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("10"),col=colors[1]) # at exprime la position en x

# print("5")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS_50$"maxlatency", xaxt="n",at=c(50),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("20"),col=colors[1]) # at exprime la position en x
  

# legend("topleft","IM-CRDT",col=colors[1], pch=1)
# legend("topright","IPFS Alone",col=colors[2], pch=1)

# print("6")
# #print(CSVFrameCRDT_IPFS_100$"maxlatency")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameIPFS_ALONE_2$"maxlatency", xaxt="n",at=c(3),add = TRUE, boxwex = 1, outlier.shape = NA, data=mtcars, names=c("1"),col=colors[2]) # at exprime la position en x

# print("7")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameIPFS_ALONE_5$"maxlatency", xaxt="n",at=c(6),add = TRUE,boxwex = 1, data=mtcars, names=c("2"),col=colors[2]) # at exprime la position en x

# print("8")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameIPFS_ALONE_10$"maxlatency", xaxt="n",at=c(11),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("5"),col=colors[2]) # at exprime la position en x

# print("9")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameIPFS_ALONE_20$"maxlatency", xaxt="n",at=c(21),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("10"),col=colors[2]) # at exprime la position en x

# print("10")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameIPFS_ALONE_50$"maxlatency", xaxt="n",at=c(51),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("20"),col=colors[2]) # at exprime la position en x
  


# par(new=TRUE)
# gna2 = boxplot(CSVFrame2$"maxlatency", xaxt="n",at=c(2),add = TRUE,boxwex = 1,data=mtcars, names=c("IFPS_alone"),col=colors[2]) # at exprime la position en x


# dev.off()


#///// STOP OF USED PART

# pdf("Output2.pdf")
# # CSVFrame1=rbind(CSVFrame1,c(1,1))
# boxplot(c(CSVFrameIPFS_ALONE100, 1,1,1, 1, 1),ylim=c(0,600),outline=FALSE,at= c(1,2,3, 4,5,6,7),outlier.shape = NA,  names=c("100","1000","10000", "", "100", "1000", "10000"), xlab="Number of Updates",ylab="Maximum latency per update (ms)", boxfill = NA, border = NA) #invisible boxes - only axes and plot area

# #print(CSVFrameCRDT_IPFS100$"maxlatency")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS100$"maxlatency", xaxt="n",at=c(1),add = TRUE, boxwex = 1, outlier.shape = NA, data=mtcars, names=c("100"),col=colors[1]) # at exprime la position en x

# print("2")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS1000$"maxlatency", xaxt="n",at=c(2),add = TRUE,boxwex = 1, data=mtcars, names=c("1000"),col=colors[1]) # at exprime la position en x

# print("3")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS10000$"maxlatency", xaxt="n",at=c(3),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("10000"),col=colors[1]) # at exprime la position en x

# print("4")

# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS100_CONCU$"maxlatency", xaxt="n",at=c(5),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("100"),col=colors[3]) # at exprime la position en x

# print("5")
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS1000_CONCU$"maxlatency", xaxt="n",at=c(6),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("1000"),col=colors[3]) # at exprime la position en x
  
# par(new=TRUE)
# gna1 = boxplot(CSVFrameCRDT_IPFS10000_CONCU$"maxlatency", xaxt="n",at=c(7),add = TRUE,boxwex = 1, outlier.shape = NA, data=mtcars, names=c("10000"),col=colors[3]) # at exprime la position en x

# legend("bottomleft",        "CRDT+IPFS",col=colors[1], pch=1)
# legend("bottomright",        "CRDT+IPFS 4 writers",col=colors[3], pch=1)

# # par(new=TRUE)
# # gna2 = boxplot(CSVFrame2$"maxlatency", xaxt="n",at=c(2),add = TRUE,boxwex = 1,data=mtcars, names=c("IFPS_alone"),col=colors[2]) # at exprime la position en x


# dev.off()

















# data_Frame = read.csv("DATA_Experience/IM-CRDT/output_2Peers_1Updater_10Updates.txt")
# data_Frame <- data_Frame %>%
#   add_column(numUpdates = 10) %>%
#   add_column(numberPeers = 2) %>%
#   add_column(numberPeerUpdating = 1) %>%
#   add_column(System = "IM-CRDT")

# a = read.csv("DATA_Experience/IM-CRDT/output_2Peers_1Updater_100Updates.csv")
# b <- data_Frame %>%
#   add_column(numUpdates = 100) %>%
#   add_column(numberPeer = 2) %>%
#   add_column(numberUpdater = 1) %>%
#   add_column(System = "IM-CRDT")

# CSVFrameIM_CRDT_2_1_1000 = read.csv("DATA_Experience/IM-CRDT/output_2Peers_1Updater_1000Updates.csv")

# CSVFrameIM_CRDT_2_2_10 = read.csv("DATA_Experience/IM-CRDT/output_2Peers_2Updater_10Updates.csv")
# CSVFrameIM_CRDT_2_2_100 = read.csv("DATA_Experience/IM-CRDT/output_2Peers_2Updater_100Updates.csv")
# CSVFrameIM_CRDT_2_2_1000 = read.csv("DATA_Experience/IM-CRDT/output_2Peers_2Updater_1000Updates.csv")


# # 5 Peers
# CSVFrameIM_CRDT_5_1_10 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_1Updater_10Updates.csv")
# CSVFrameIM_CRDT_5_1_100 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_1Updater_100Updates.csv")
# CSVFrameIM_CRDT_5_1_1000 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_1Updater_1000Updates.csv")

# CSVFrameIM_CRDT_5_2_10 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_2Updater_10Updates.csv")
# CSVFrameIM_CRDT_5_2_100 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_2Updater_100Updates.csv")
# CSVFrameIM_CRDT_5_2_1000 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_2Updater_1000Updates.csv")

# CSVFrameIM_CRDT_5_5_10 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_5Updater_10Updates.csv")
# CSVFrameIM_CRDT_5_5_100 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_5Updater_100Updates.csv")
# CSVFrameIM_CRDT_5_5_1000 = read.csv("DATA_Experience/IM-CRDT/output_5Peers_5Updater_1000Updates.csv")


# # 10 Peers
# CSVFrameIM_CRDT_10_1_10 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_1Updater_10Updates.csv")
# CSVFrameIM_CRDT_10_1_100 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_1Updater_100Updates.csv")
# CSVFrameIM_CRDT_10_1_1000 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_1Updater_1000Updates.csv")

# CSVFrameIM_CRDT_10_2_10 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_2Updater_10Updates.csv")
# CSVFrameIM_CRDT_10_2_100 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_2Updater_100Updates.csv")
# CSVFrameIM_CRDT_10_2_1000 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_2Updater_1000Updates.csv")

# CSVFrameIM_CRDT_10_5_10 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_5Updater_10Updates.csv")
# CSVFrameIM_CRDT_10_5_100 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_5Updater_100Updates.csv")
# CSVFrameIM_CRDT_10_5_1000 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_5Updater_1000Updates.csv")

# CSVFrameIM_CRDT_10_10_10 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_10Updater_10Updates.csv")
# CSVFrameIM_CRDT_10_10_100 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_10Updater_100Updates.csv")
# CSVFrameIM_CRDT_10_10_1000 = read.csv("DATA_Experience/IM-CRDT/output_10Peers_10Updater_1000Updates.csv")


# # 20 Peers
# CSVFrameIM_CRDT_20_1_10 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_1Updater_10Updates.csv")
# CSVFrameIM_CRDT_20_1_100 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_1Updater_100Updates.csv")
# CSVFrameIM_CRDT_20_1_1000 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_1Updater_1000Updates.csv")

# CSVFrameIM_CRDT_20_2_10 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_2Updater_10Updates.csv")
# CSVFrameIM_CRDT_20_2_100 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_2Updater_100Updates.csv")
# CSVFrameIM_CRDT_20_2_1000 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_2Updater_1000Updates.csv")

# CSVFrameIM_CRDT_20_5_10 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_5Updater_10Updates.csv")
# CSVFrameIM_CRDT_20_5_100 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_5Updater_100Updates.csv")
# CSVFrameIM_CRDT_20_5_1000 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_5Updater_1000Updates.csv")

# CSVFrameIM_CRDT_20_10_10 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_10Updater_10Updates.csv")
# CSVFrameIM_CRDT_20_10_100 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_10Updater_100Updates.csv")
# CSVFrameIM_CRDT_20_10_1000 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_10Updater_1000Updates.csv")

# CSVFrameIM_CRDT_20_20_10 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_20Updater_10Updates.csv")
# CSVFrameIM_CRDT_20_20_100 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_20Updater_100Updates.csv")
# CSVFrameIM_CRDT_20_20_1000 = read.csv("DATA_Experience/IM-CRDT/output_20Peers_20Updater_1000Updates.csv")


# # 50 Peers
# CSVFrameIM_CRDT_50_1_10 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_1Updater_10Updates.csv")
# CSVFrameIM_CRDT_50_1_100 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_1Updater_100Updates.csv")
# CSVFrameIM_CRDT_50_1_1000 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_1Updater_1000Updates.csv")

# CSVFrameIM_CRDT_50_5_10 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_2Updater_10Updates.csv")
# CSVFrameIM_CRDT_50_5_100 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_2Updater_100Updates.csv")
# CSVFrameIM_CRDT_50_5_1000 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_2Updater_1000Updates.csv")

# CSVFrameIM_CRDT_50_5_10 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_5Updater_10Updates.csv")
# CSVFrameIM_CRDT_50_5_100 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_5Updater_100Updates.csv")
# CSVFrameIM_CRDT_50_5_1000 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_5Updater_1000Updates.csv")

# CSVFrameIM_CRDT_50_10_10 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_10Updater_10Updates.csv")
# CSVFrameIM_CRDT_50_10_100 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_10Updater_100Updates.csv")
# CSVFrameIM_CRDT_50_10_1000 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_10Updater_1000Updates.csv")

# CSVFrameIM_CRDT_50_20_10 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_20Updater_10Updates.csv")
# CSVFrameIM_CRDT_50_20_100 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_20Updater_100Updates.csv")
# CSVFrameIM_CRDT_50_20_1000 = read.csv("DATA_Experience/IM-CRDT/output_50Peers_20Updater_1000Updates.csv")

# #IPFS Alone
# # 2 Peers
# CSVFrameIPFS_Alone_2_1_10 = read.csv("DATA_Experience/IPFS_Alone/output_2Peers_1Updater_10Updates.csv")
# CSVFrameIPFS_Alone_2_1_100 = read.csv("DATA_Experience/IPFS_Alone/output_2Peers_1Updater_100Updates.csv")
# CSVFrameIPFS_Alone_2_1_1000 = read.csv("DATA_Experience/IPFS_Alone/output_2Peers_1Updater_1000Updates.csv")

# # 5 Peers
# CSVFrameIPFS_Alone_5_1_10 = read.csv("DATA_Experience/IPFS_Alone/output_5Peers_1Updater_10Updates.csv")
# CSVFrameIPFS_Alone_5_1_100= read.csv("DATA_Experience/IPFS_Alone/output_5Peers_1Updater_100Updates.csv")
# CSVFrameIPFS_Alone_5_1_1000 = read.csv("DATA_Experience/IPFS_Alone/output_5Peers_1Updater_1000Updates.csv")


# # 10 Peers
# CSVFrameIPFS_Alone_10_1_10 = read.csv("DATA_Experience/IPFS_Alone/output_10Peers_1Updater_10Updates.csv")
# CSVFrameIPFS_Alone_10_1_100 = read.csv("DATA_Experience/IPFS_Alone/output_10Peers_1Updater_100Updates.csv")
# CSVFrameIPFS_Alone_10_1_1000 = read.csv("DATA_Experience/IPFS_Alone/output_10Peers_1Updater_1000Updates.csv")


# # 20 Peers
# CSVFrameIPFS_Alone_20_1_10 = read.csv("DATA_Experience/IPFS_Alone/output_20Peers_1Updater_10Updates.csv")
# CSVFrameIPFS_Alone_20_1_100 = read.csv("DATA_Experience/IPFS_Alone/output_20Peers_1Updater_100Updates.csv")
# CSVFrameIPFS_Alone_20_1_1000 = read.csv("DATA_Experience/IPFS_Alone/output_20Peers_1Updater_1000Updates.csv")



# # 50 Peers
# CSVFrameIPFS_Alone_50_1_10 = read.csv("DATA_Experience/IPFS_Alone/output_50Peers_1Updater_10Updates.csv")
# CSVFrameIPFS_Alone_50_1_100 = read.csv("DATA_Experience/IPFS_Alone/output_50Peers_1Updater_100Updates.csv")
# CSVFrameIPFS_Alone_50_1_1000 = read.csv("DATA_Experience/IPFS_Alone/output_50Peers_1Updater_1000Updates.csv")


