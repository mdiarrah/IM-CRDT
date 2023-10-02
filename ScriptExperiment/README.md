# ScriptExperiment

These scripts are used to run he experiments, it must be pushed remotely (using `sendscript.sh G5K_SITE.g5k`)

Then, `run_multiple.sh` is the enter point, it does call recursively all the other program.

In the remote site, you will need to run`./run_multiple.sh NDOEFILE` where NODEFILE does contain all the accessible nodes.
__Notice than the bootstrapPeer will be the first in the list, then all the other nodes will be taken from the tail.__