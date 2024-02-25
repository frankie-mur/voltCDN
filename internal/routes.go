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

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", Health)
	mux.HandleFunc("POST /photo", app.CreatePhoto)
	mux.HandleFunc("GET /photo", app.GetAllPhotos)
	mux.HandleFunc("GET /photo/{id}", app.GetPhotoById)
	mux.HandleFunc("DELETE /photo/{id}", app.DeletePhotoById)

	mux.HandleFunc("GET /", app.IndexPage)
	return mux
}

type PageData struct {
	Body   string
	Photos []*db.PhotoEntity
}

func (app *Application) IndexPage(w http.ResponseWriter, r *http.Request) {
	indexPage := "./frontend/pages/index.tmpl"
	data, err := app.db.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := PageData{
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

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

func (app *Application) CreatePhoto(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) GetAllPhotos(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) GetPhotoById(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) DeletePhotoById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	if idString == "" {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}

	err := app.db.Delete(idString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Here we just write an empty response, as this is an HTMX request
	//And the response will replace the element specified in the request.
	w.Write([]byte(""))
}
