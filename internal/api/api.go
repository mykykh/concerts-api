package api

import (
    "os"
    "fmt"
    "context"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Api struct {
    db *pgxpool.Pool
    router chi.Router
}

func connectToDB() (*pgxpool.Pool) {
    dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
            os.Exit(1)
    }
    defer dbpool.Close()

    return dbpool
}

func initRouter() chi.Router {
    router := chi.NewRouter()

    router.Mount("/concerts", ConcertsResource{}.Routes())

    return router
}

func Init() Api {
    api := Api{db: connectToDB(), router: initRouter()}

    return api
}

func (api Api) Run() {
    err := http.ListenAndServe(":8080", api.router)

    if err != nil {
        fmt.Fprint(os.Stderr, err)
    }
}
