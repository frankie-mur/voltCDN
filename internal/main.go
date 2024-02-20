package main

import (
	"net/http"

	"github.com/frankie-mur/voltCDN/internal/db"
)

type Application struct {
	db *db.SQLLiteRepository
}

func main() {
	app := Application{
		db: db.NewSQLLiteRepository(),
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /entities", app.createEntity)
	mux.HandleFunc("GET /entities/{id}", app.getEntityById)

	http.ListenAndServe(":8080", mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

func (app *Application) createEntity(w http.ResponseWriter, r *http.Request) {
	e := &db.Entity{
		Id:   "1",
		Name: "test",
	}
	err := app.db.Save(*e)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte("created"))
}

func (app *Application) getEntityById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	res, err := app.db.Get(idString)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(res.Name))
}
