package mapreduce

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"os"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// TODO:
	// You will need to write this function.
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	//
	// The intermediate output of a map task is stored in the file
	// system as multiple files whose name indicates which map task produced
	// them, as well as which reduce task they are for. Coming up with a
	// scheme for how to store the key/value pairs on disk can be tricky,
	// especially when taking into account that both keys and values could
	// contain newlines, quotes, and any other character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	// Use checkError to handle errors.

	// 1. Read the file
	data, err := ioutil.ReadFile(inFile)
	checkError(err)

	// get the contents of the file
	fcontent := string(data)
	// get the key value pairs
	keyVal := mapF(inFile, fcontent)

	// loop trough nreduce to create that number of files
	for i := 0; i < nReduce; i++ {
		//get the name of the temp file
		tmpFileName := reduceName(jobName, mapTaskNumber, i)
		// create a new tmp file with the reduced name
		tmpf, err := os.OpenFile(tmpFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkError(err)
		//frite encoded data to the file
		enc := json.NewEncoder(tmpf)
		for _, kv := range keyVal {

			//find if a key belongs in the current file
			hashedReduced := ihash(kv.Key) % uint32(nReduce)
			if hashedReduced == uint32(i) {
				err := enc.Encode(&kv)
				checkError(err)
			}
		}
		//close the file when done
		err = tmpf.Close()
		checkError(err)
	}

}

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
