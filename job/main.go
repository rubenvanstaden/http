package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Task struct {
	Type    string `json:"type"`
	Queue   string `json:"queue"`
	Payload []byte `json:"payload"`
}

// In-memory mock of task queue.
type Tasks []*Task

func (s Tasks) Enqueue(task *Task) error {
	s = append(s, task)
	return nil
}

// Job domain event.
type JobEvent struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

type handler struct {
	tasks *Tasks
}

func (s handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:

		var event JobEvent

		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Can't decode body", http.StatusBadRequest)
			return
		}

		bytes, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		task := &Task{
			Type:    event.Status,
			Queue:   "slurm",
			Payload: bytes,
		}

		err = s.tasks.Enqueue(task)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(s.tasks)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	tasks := Tasks{}

	mux := http.NewServeMux()
	mux.Handle("/update", handler{&tasks})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
