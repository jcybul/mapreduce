package mapreduce

import "sync"

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO:

	//loop trough n tasks and assign them
	var wg sync.WaitGroup
	println(ntasks)
	for i := 0; i < ntasks; i++ {

		// add new task to wait group
		wg.Add(1)
		// create dtask args
		dTask := new(DoTaskArgs)
		dTask.JobName = mr.jobName
		dTask.File = mr.files[i]
		dTask.Phase = phase
		dTask.TaskNumber = i
		dTask.NumOtherPhase = nios

		go func() {
			for { // loop in case of failure
				w := <-mr.registerChannel                     // get open worker
				var reply ShutdownReply                       // can be any pointer but this for simplicity
				ok := call(w, "Worker.DoTask", dTask, &reply) // call do task on worker
				if ok {                                       // if no errors
					go func() { // free worker
						// in goroutine to prevent deadlock
						mr.registerChannel <- w
					}()
					// deccrement wait group
					wg.Done()
					break

				}

			}

		}()
	}
	// wait for all to finish
	wg.Wait()

	debug("Schedule: %v phase done\n", phase)
}
