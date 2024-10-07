package main

import (
	"log"

	"github.com/aranw/freq-demo/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, workflow.FrequencyBatchTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(workflow.ProcessFrequencyBatch)
	w.RegisterActivity(workflow.Min)
	w.RegisterActivity(workflow.Max)
	w.RegisterActivity(workflow.Avg)
	w.RegisterActivity(workflow.StdDev)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
