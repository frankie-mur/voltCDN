package main

import (
	"encoding/json"
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
	mux.HandleFunc("POST /photo", app.createEntity)
	mux.HandleFunc("GET /photo", app.getAllEntities)
	mux.HandleFunc("GET /photo/{id}", app.getEntityById)

	http.ListenAndServe(":8080", mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

type CreateRequest struct {
	Data string `json:"data"`
}

func (app *Application) createEntity(w http.ResponseWriter, r *http.Request) {
	jsonRequest := r.Body
	defer jsonRequest.Close()

	var request CreateRequest
	err := json.NewDecoder(jsonRequest).Decode(&request)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	if request.Data == "" {
		http.Error(w, "data is required", http.StatusBadRequest)
		return
	}

	e := &db.PhotoEntity{
		Data: request.Data,
	}
	err = app.db.Save(*e)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte("created"))
}

func (app *Application) getAllEntities(w http.ResponseWriter, r *http.Request) {
	entities, err := app.db.GetAll()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	resp, err := json.Marshal(entities)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(resp)
}

func (app *Application) getEntityById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	res, err := app.db.Get(idString)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(res.Data))
}
