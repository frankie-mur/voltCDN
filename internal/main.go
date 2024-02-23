package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

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

	mux.HandleFunc("GET /", app.indexPage)

	http.ListenAndServe(":8080", mux)
}

type PageData struct {
	Body   string
	Photos []*db.PhotoEntity
}

func (app *Application) indexPage(w http.ResponseWriter, r *http.Request) {
	indexPage := "./frontend/pages/index.tmpl"
	data, err := app.db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Body:   "Hello World",
		Photos: data,
	}

	tmpl, err := template.ParseFiles(indexPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

func (app *Application) createPhoto(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) //10mb
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	file, handler, err := r.FormFile("img")

	if err != nil {
		http.Error(w, "img is required", http.StatusBadRequest)
		return
	}
	defer file.Close()
	// Create a temporary file to store the uploaded image
	tempFile, err := os.CreateTemp("", "temp-image-*.png")
	if err != nil {
		http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	// Read the content of the temporary file
	imageData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	e := &db.PhotoEntity{
		Name: handler.Filename,
		Data: imageData,
	}
	fmt.Printf("Saving %s\n", e.Name)
	err = app.db.Save(*e)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Photo saved"))
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
	fmt.Println(res.Name)
	contentType := filepath.Ext(res.Name)[1:]
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", contentType))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(res.Data)))

	if _, err := w.Write(res.Data); err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to write image to response", http.StatusInternalServerError)
		return
	}
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
}
