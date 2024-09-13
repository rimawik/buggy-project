package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB



func main() {

	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	if err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}

	defer db.Close()

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		rows, err := db.Query("SELECT name FROM users")
		if err != nil {
			log.Printf("couldn't get users: %v", err.Error())
			http.Error(w, "database_error", http.StatusInternalServerError)
			return

		}
		defer rows.Close()

		for rows.Next() {
			var name string
			rows.Scan(&name)
			fmt.Fprintf(w, "User: %s\n", name)
		}
	}()

	wg.Wait()
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		time.Sleep(5 * time.Second) // Simulate a long database operation

		username := r.URL.Query().Get("name")
		_, err := db.Exec("INSERT INTO users (name) VALUES ('" + username + "')")

		if err != nil {
			log.Printf("Failed to create user:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to create user: %v", err)
			return
		}

		fmt.Fprintf(w, "User %s created successfully", username)
	}()

	wg.Wait()
}
