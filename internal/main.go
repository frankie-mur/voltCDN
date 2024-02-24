package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/frankie-mur/voltCDN/internal/db"
)

type Application struct {
	db *db.SQLLiteRepository
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Port to listen on")

	flag.Parse()
	app := Application{
		db: db.NewSQLLiteRepository(),
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	fmt.Printf("Listening on port %s\n", port)
	srv.ListenAndServe()
}
