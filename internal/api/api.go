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

    return dbpool
}

func initRouter(db *pgxpool.Pool) chi.Router {
    router := chi.NewRouter()

    router.Mount("/concerts", ConcertsResource{db: db}.Routes())

    return router
}

func Init() Api {
    db := connectToDB()
    router := initRouter(db)

    return Api{db: db, router: router};
}

func (api Api) Run() {
    err := http.ListenAndServe(":8080", api.router)

    if err != nil {
        fmt.Fprint(os.Stderr, err)
    }
}
