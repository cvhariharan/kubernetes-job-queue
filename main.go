package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cvhariharan/kubernetes-job-queue/jobqueue"
	"github.com/google/uuid"
)

func main() {
	log.SetFlags(log.Lshortfile)
	var worker, exitIfEmpty bool
	flag.BoolVar(&worker, "WORKER", false, "Worker node")
	flag.BoolVar(&exitIfEmpty, "EXIT", false, "Exit if the queue is empty")
	flag.Parse()

	queue := jobqueue.NewRedisJobqueue()

	if worker {
		host, _ := os.Hostname()
		log.Println("Starting worker", host)
		jobsChan := queue.Subscribe(context.Background(), exitIfEmpty)
		for {
			job, ok := <-jobsChan
			if !ok {
				return
			}
			time.Sleep(5 * time.Second)
			log.Println(fmt.Sprintf("Processed job %s by worker %s", job, host))
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pubJobqueue := make(chan string)
	go queue.Publish(context.Background(), pubJobqueue)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pubJobqueue <- uuid.New().String()
		w.WriteHeader(http.StatusAccepted)
	})

	log.Println("Starting API server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
