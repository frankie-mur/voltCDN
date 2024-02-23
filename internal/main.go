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
	mux.HandleFunc("POST /photo", app.createPhoto)
	mux.HandleFunc("GET /photo", app.getAllPhotos)
	mux.HandleFunc("GET /photo/{id}", app.getPhotoById)
	mux.HandleFunc("DELETE /photo/{id}", app.deletePhotoById)

	http.ListenAndServe(":8080", mux)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

type CreateRequest struct {
	Data string `json:"data"`
}

func (app *Application) createPhoto(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusCreated)
	return
}

func (app *Application) getAllPhotos(w http.ResponseWriter, r *http.Request) {
	entities, err := app.db.GetAll()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	resp, err := json.Marshal(entities)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(resp)
	return
}

func (app *Application) getPhotoById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}

	res, err := app.db.Get(idString)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(jsonData)
}

func (app *Application) deletePhotoById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}

	err := app.db.Delete(idString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
