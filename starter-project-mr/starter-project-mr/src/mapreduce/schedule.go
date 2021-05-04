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
	println("in schedule")
	//loop trough n tasks and assign them
	var wg sync.WaitGroup
	println(ntasks)
	for i := 0; i < ntasks; i++ {

		println(i)
		wg.Add(1)
		// w := new(Worker)
		dTask := new(DoTaskArgs)
		dTask.JobName = mr.jobName
		dTask.File = mr.files[i]
		dTask.Phase = phase
		dTask.TaskNumber = i
		dTask.NumOtherPhase = nios
		go func() {
			for {
				w := <-mr.registerChannel
				var reply ShutdownReply
				ok := call(w, "Worker.DoTask", dTask, &reply)
				if ok {
					go func() {
						mr.registerChannel <- w
					}()
					wg.Done()
					break

				}

			}

		}()
	}
	println("waiting")
	wg.Wait()

	debug("Schedule: %v phase done\n", phase)
}
