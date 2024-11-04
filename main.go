package main

import (
    "net/http"
    "context"
    "fmt"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mykykh/concerts-api/internal/api"
)

func testDBConnect() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

func main() {
    r := chi.NewRouter()

    r.Mount("/concerts", api.ConcertsResource{}.Routes())

    testDBConnect()

    http.ListenAndServe(":8080", r)
}
