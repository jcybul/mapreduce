package mapreduce

import (
	"encoding/json"
	"os"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Use checkError to handle errors.

	//1 loop trough the files we created in map

	kvmap := make(map[string][]string)

	for i := 0; i < nMap; i++ {

		tmpTaskNumber := reduceName(jobName, i, reduceTaskNumber)

		tmpf, err := os.OpenFile(tmpTaskNumber, os.O_RDONLY, 0644)
		checkError(err)

		dec := json.NewDecoder(tmpf)

		for dec.More() {
			emptyStruct := new(KeyValue)

			err := dec.Decode(&emptyStruct)
			checkError(err)

			tmpList := kvmap[emptyStruct.Key]
			kvmap[emptyStruct.Key] = append(tmpList, emptyStruct.Value)

		}

		tmpf.Close()

	}

	mn := mergeName(jobName, reduceTaskNumber)

	openmn, err := os.OpenFile(mn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)

	enc := json.NewEncoder(openmn)
	for key, val := range kvmap {
		wrt := reduceF(key, val)
		enc.Encode(KeyValue{key, wrt})
	}
	openmn.Close()

}
