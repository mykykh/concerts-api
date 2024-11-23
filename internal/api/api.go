package api

import (
    "os"
    "fmt"
    "context"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/swaggo/http-swagger/v2"

    _ "github.com/mykykh/concerts-api/docs"
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

// @title Concerts api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description To obtain token go to [/auth/login](/auth/login). You also need to add Bearer before pasting it belove. It should look like: Bearer your-access-token
func initRouter(db *pgxpool.Pool) chi.Router {
    router := chi.NewRouter()

    router.Mount("/auth", AuthResource{db: db}.Routes())
    router.Mount("/users", UsersResource{db: db}.Routes())
    router.Mount("/concerts", ConcertsResource{db: db}.Routes())

    router.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"),
    ))

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
