CS4513: Project 3 The MapReduce Library
=======================================

Note, this document includes a number of design questions that can help your implementation. We highly recommend that you answer each design question **before** attempting the corresponding implementation.
These questions will help you design and plan your implementation and guide you towards the resources you need.
Finally, if you are unsure how to start the project, we recommend you visit office hours for some guidance on these questions before attempting to implement this project.


Team members
-----------------

1. Truman Larson (tlarson@wpi.edu)
2. Joseph Cybul (jcybul@wpi.edu)

Design Questions
------------------

(2 point) 1. If there are n input files, and nReduce number of reduce tasks , how does the the MapReduce Library uniquely name the intermediate files?

f0-0, ... f0-[nReduce-1]
f1-0, ... f1-[nReduce-1]
...
f[n-1]-0 ... f[n-1][nReduce-1]


(1 point) 2. Following the previous question, for the reduce task r, what are the names of files will it work on?

f0-r,
f1-r,
...
f[n-1]-r


(1 point) 3. If the submitted mapreduce job name is "test", what will be the final output file's name?

mrtmp.test


(2 point) 4. Based on `mapreduce/test_test.go`, when you run the `TestBasic()` function, how many master and workers will be started? And what are their respective addresses and their naming schemes?


1 master, 2 workers.
Addresses should look similar to:

/var/tmp/824-UID/mrPID-master
/var/tmp/824-UID/mrPID-worker0
/var/tmp/824-UID/mrPID-worker1

UID and PID can be any integers.


(4 point) 5. In real-world deployments, when giving a mapreduce job, we often start master and workers on different machines (physical or virtual). Describe briefly the protocol that allows master and workers be aware of each other's existence, and subsequently start working together on completing the mapreduce job. Your description should be grounded on the RPC communications.


(Most students probably will be able to answer this question correctly---as the process is summarized in the project description. The purpose of this question is to get students to read and understand the protocol.)

Just give submission full points for this question.


(2 point) 6. The current design and implementation uses a number of RPC methods. Can you find out all the RPCs and list their signatures? Briefly describe the criteria
a method needs to satisfy to be considered a RPC method. (Hint: you can look up at: https://golang.org/pkg/net/rpc/)

Worker.DoTask
Worker.Shutdown
Master.Shutdown
Master.Register

The method and its type must be exported. 

It must have two types with the second type being a pointer to a reply block. 

It must return and error type.

7. We provide a bash script called `main/test-wc-distributed.sh` that allows testing the distributed MapReduce implementation.

Errata
------

Describe any known errors, bugs, or deviations from the requirements.

Mapping: 
    No known issues with mapping. For common_map we utilized the json encoding to write to all of the
reduce files. To determine if a given key should be in a corresponding file, we use the hash mod the
number of reduce files and include it if that matches the current file number.

Reducing:
    No known errors. Two core loops control this functionality: 1 to get the mapped contents into memory
and another to call the reduce function on each key value pair.

Scheduling:
No known errors.

    For each ntask, we generate a dtask args struct and start a goroutine to call a worker. We add 
1 to the wait group to indicate that we have another task to complete.To find an avaiable worker,
we get one from the mr.registerChannel and call Worker.DoTask on that. If the result is without 
error, we try to free that worker by putting it back in the mr.registerChannel and decrement 
the wait group. If there is an error with that worker, we will loop back and try to get a new 
worker and repeat the process. After all workers are scheduled, we wait for the wait group to 
finish which indicates that we are done with this phase. 

---
