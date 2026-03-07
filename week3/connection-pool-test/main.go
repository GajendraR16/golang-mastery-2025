package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	testNonPooled()

	testPooled()
}

func testNonPooled() {
	connStr := "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()

	//No Configuration - uses default

	start := time.Now()
	runQueries(db, 100)
	fmt.Printf("Non Pooled: %v\n", time.Since(start))
}

func testPooled() {
	connStr := "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=diasble"
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()

	//Configuration

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	start := time.Now()
	runQueries(db, 100)
	fmt.Printf("Pooled : %v\n", time.Since(start))
}

func runQueries(db *sql.DB, count int) {
	var wg sync.WaitGroup
	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var count int
			db.QueryRow("Select Count(*) from tasks").Scan(&count)
		}()
	}
	wg.Wait()
}
