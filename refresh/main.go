package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Result struct {
	URL    string `json:"url"`
	Status int    `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Request struct {
	URLs []string `json:"urls"`
}

func check(ctx context.Context, urls []string) []Result {
	jobs := make(chan string)
	out := make(chan Result)
	client := &http.Client{}

	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					out <- Result{URL: url, Error: err.Error()}
					continue
				}
				resp, err := client.Do(req)
				if err != nil {
					out <- Result{URL: url, Error: err.Error()}
					continue
				}
				out <- Result{URL: url, Status: resp.StatusCode}
				defer resp.Body.Close()
			}
		}()
	}

	go func() {
		for _, u := range urls {
			jobs <- u
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	var res []Result
	for r := range out {
		res = append(res, r)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Status < res[j].Status
	})

	return res
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	results := check(ctx, req.URLs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/check", checkHandler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
